package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/hashicorp/vault/api"
)

type VaultStore struct {
	client *api.Client
	kvPath string
}

type Config struct {
	MSPID         string
	Org           string
	TLSCertPath   string
	IngressDomain string
	Namespace     string
	VaultAddr     string
	VaultToken    string
	KVPath        string
	RCAMSPPath    string
}

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	UserID   string `json:"user_id" binding:"required"`
}

func NewVaultStore(cfg *Config) (*VaultStore, error) {
	config := api.DefaultConfig()
	config.Address = cfg.VaultAddr
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	client.SetToken(cfg.VaultToken)

	return &VaultStore{
		client: client,
		kvPath: cfg.KVPath,
	}, nil
}

func loadConfig() *Config {
	return &Config{
		MSPID:         getEnvWithDefault("MSP_ID", "Org1MSP"),
		Org:           getEnvWithDefault("ORG", "org1"),
		TLSCertPath:   getEnvWithDefault("TLS_CERT_PATH", "/etc/tls/tls-cert.pem"),
		IngressDomain: getEnvWithDefault("WORKSHOP_INGRESS_DOMAIN", "localho.st"),
		Namespace:     getEnvWithDefault("WORKSHOP_NAMESPACE", "test-network"),
		VaultAddr:     getEnvWithDefault("VAULT_ADDR", "http://127.0.0.1:8200"),
		VaultToken:    getEnvWithDefault("VAULT_TOKEN", "<vault-token>"),
		KVPath:        getEnvWithDefault("VAULT_KV_PATH", "fabric/msp"),
		RCAMSPPath:    getEnvWithDefault("RCAMSP_PATH", "/etc/rcaadmin/msp"),
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	cfg := loadConfig()
	vault, err := NewVaultStore(cfg)
	if err != nil {
		log.Fatal("Failed to initialize Vault client: ", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		handleRegister(w, r, cfg, vault)
	})
	r.Post("/enroll", func(w http.ResponseWriter, r *http.Request) {
		handleEnroll(w, r, cfg, vault)
	})
	r.Post("/register-enroll", func(w http.ResponseWriter, r *http.Request) {
		handleRegisterAndEnroll(w, r, cfg, vault)
	})
	r.Delete("/revoke/{username}", func(w http.ResponseWriter, r *http.Request) {
		handleRevoke(w, r, cfg, vault)
	})
	r.Get("/msp/{username}", func(w http.ResponseWriter, r *http.Request) {
		handleGetMSP(w, r, cfg, vault)
	})

	log.Println("Server started on :8080")
	if err = http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request, cfg *Config, vault *VaultStore) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderError(w, r, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.UserID == "" {
		renderError(w, r, "Username and UserID are required", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		password, err := generatePassword()
		if err != nil {
			renderError(w, r, "Failed to generate password", http.StatusInternalServerError)
			return
		}
		req.Password = password
	}

	err := registerUser(cfg, req.Username, req.Password, req.UserID)
	if err != nil {
		renderError(w, r, fmt.Sprintf("Registration failed: %v", err), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"username": req.Username,
		"password": req.Password,
		"message":  "User registered successfully",
	})
}

func handleEnroll(w http.ResponseWriter, r *http.Request, cfg *Config, vault *VaultStore) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderError(w, r, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		renderError(w, r, "Username and password are required", http.StatusBadRequest)
		return
	}

	mspPath, err := enrollUser(cfg, req.Username, req.Password)
	if err != nil {
		renderError(w, r, fmt.Sprintf("Enrollment failed: %v", err), http.StatusInternalServerError)
		return
	}

	if err := vault.StoreUserMSP(req.Username, mspPath); err != nil {
		renderError(w, r, fmt.Sprintf("Vault storage failed: %v", err), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"username": req.Username,
		"message":  "User enrolled and MSP stored successfully",
	})
}

func handleRegisterAndEnroll(w http.ResponseWriter, r *http.Request, cfg *Config, vault *VaultStore) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderError(w, r, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		password, err := generatePassword()
		if err != nil {
			renderError(w, r, "Failed to generate password", http.StatusInternalServerError)
			return
		}
		req.Password = password
	}

	if err := registerUser(cfg, req.Username, req.Password, req.UserID); err != nil {
		renderError(w, r, fmt.Sprintf("Registration failed: %v", err), http.StatusInternalServerError)
		return
	}

	mspPath, err := enrollUser(cfg, req.Username, req.Password)
	if err != nil {
		_ = revokeUser(cfg, vault, req.Username, "enrollment failed")
		renderError(w, r, fmt.Sprintf("Enrollment failed: %v", err), http.StatusInternalServerError)
		return
	}

	if err := vault.StoreUserMSP(req.Username, mspPath); err != nil {
		_ = revokeUser(cfg, vault, req.Username, "vault storage failed")
		renderError(w, r, fmt.Sprintf("Vault storage failed: %v", err), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"username": req.Username,
		"password": req.Password,
		"message":  "User registered, enrolled and MSP stored successfully",
	})
}

func handleRevoke(w http.ResponseWriter, r *http.Request, cfg *Config, vault *VaultStore) {
	username := chi.URLParam(r, "username")
	if username == "" {
		renderError(w, r, "Username is required", http.StatusBadRequest)
		return
	}

	if err := revokeUser(cfg, vault, username, "admin revocation"); err != nil {
		renderError(w, r, fmt.Sprintf("Revocation failed: %v", err), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"username": username,
		"message":  "User revoked and MSP removed successfully",
	})
}

func renderError(w http.ResponseWriter, r *http.Request, msg string, status int) {
	render.Status(r, status)
	render.JSON(w, r, map[string]interface{}{
		"error":   http.StatusText(status),
		"message": msg,
	})
}

func generatePassword() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func registerUser(cfg *Config, username, password, userID string) error {
	caAddress := fmt.Sprintf("%s-%s-ca-ca.%s", cfg.Namespace, cfg.Org, cfg.IngressDomain)

	cmd := exec.Command("fabric-ca-client", "register",
		"--id.name", username,
		"--id.secret", password,
		"--id.type", "client",
		"--id.affiliation", cfg.Org,
		"--id.attrs", fmt.Sprintf("identity.id=%s:ecert", userID),
		"--url", fmt.Sprintf("https://%s", caAddress),
		"--tls.certfiles", cfg.TLSCertPath,
		"--mspdir", cfg.RCAMSPPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("registration failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

func enrollUser(cfg *Config, username, password string) (string, error) {
	tempDir, err := ioutil.TempDir("", "msp-")
	if err != nil {
		return "", err
	}

	caAddress := fmt.Sprintf("%s-%s-ca-ca.%s", cfg.Namespace, cfg.Org, cfg.IngressDomain)
	enrollURL := fmt.Sprintf("https://%s:%s@%s", username, password, caAddress)

	cmd := exec.Command("fabric-ca-client", "enroll",
		"--url", enrollURL,
		"--tls.certfiles", cfg.TLSCertPath,
		"--mspdir", tempDir,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("enrollment failed: %v\nOutput: %s", err, string(output))
	}

	return tempDir, nil
}

func revokeUser(cfg *Config, vault *VaultStore, username, reason string) error {
	// Revoke using RCA admin MSP
	cmd := exec.Command("fabric-ca-client", "revoke",
		"--revoke.name", username,
		"--revoke.reason", reason,
		"--url", fmt.Sprintf("https://%s-%s-ca-ca.%s", cfg.Namespace, cfg.Org, cfg.IngressDomain),
		"--tls.certfiles", cfg.TLSCertPath,
		"--mspdir", cfg.RCAMSPPath,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("revocation failed: %v\nOutput: %s", err, string(output))
	}

	// Delete all versions and metadata
	path := fmt.Sprintf("%s/metadata/users/%s", vault.kvPath, username)
	_, err := vault.client.Logical().Delete(path)
	if err != nil {
		return fmt.Errorf("vault deletion failed: %v", err)
	}

	// Destroy all versions
	destroyPath := fmt.Sprintf("%s/destroy/users/%s", vault.kvPath, username)
	_, err = vault.client.Logical().Write(destroyPath, map[string]interface{}{
		"versions": []int{1},
	})

	return nil
}

// Add handler function
func handleGetMSP(w http.ResponseWriter, r *http.Request, cfg *Config, vault *VaultStore) {
	username := chi.URLParam(r, "username")
	if username == "" {
		renderError(w, r, "Username is required", http.StatusBadRequest)
		return
	}

	mspData, err := vault.GetUserMSP(username)
	if err != nil {
		renderError(w, r, fmt.Sprintf("Failed to retrieve MSP: %v", err), http.StatusNotFound)
		return
	}

	// Convert []byte values to base64 for JSON serialization
	response := make(map[string]string)
	for path, content := range mspData {
		response[path] = string(content)
	}

	render.JSON(w, r, map[string]interface{}{
		"username": username,
		"msp":      response,
	})
}

func (v *VaultStore) StoreUserMSP(username, mspPath string) error {
	// Store IssuerRevocationPublicKey if exists
	revocationKeyPath := filepath.Join(mspPath, "IssuerRevocationPublicKey")
	if content, err := ioutil.ReadFile(revocationKeyPath); err == nil {
		vaultPath := fmt.Sprintf("%s/data/users/%s/IssuerRevocationPublicKey", v.kvPath, username)
		if _, err := v.client.Logical().Write(vaultPath, map[string]interface{}{
			"data": map[string]interface{}{
				"content": base64.StdEncoding.EncodeToString(content),
			},
		}); err != nil {
			return fmt.Errorf("failed to store revocation key: %v", err)
		}
	}
	// Store signcerts
	signcertsDir := filepath.Join(mspPath, "signcerts")
	if err := v.storeDirectory(username, signcertsDir); err != nil {
		return fmt.Errorf("failed to store signcerts: %v", err)
	}

	// Store keystore
	keystoreDir := filepath.Join(mspPath, "keystore")
	if err := v.storeDirectory(username, keystoreDir); err != nil {
		return fmt.Errorf("failed to store keystore: %v", err)
	}

	// Store cacerts
	cacertsDir := filepath.Join(mspPath, "cacerts")
	if err := v.storeDirectory(username, cacertsDir); err != nil {
		return fmt.Errorf("failed to store cacerts: %v", err)
	}

	return nil
}

func (v *VaultStore) storeDirectory(username, dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %v", path, err)
		}

		relPath, _ := filepath.Rel(dirPath, path)
		vaultPath := fmt.Sprintf("%s/data/users/%s/%s/%s", v.kvPath, username, filepath.Base(dirPath), relPath)

		fmt.Println(vaultPath)

		_, err = v.client.Logical().Write(vaultPath, map[string]interface{}{
			"data": map[string]interface{}{
				"content": base64.StdEncoding.EncodeToString(content),
			},
		})

		return err
	})
}

func (v *VaultStore) GetUserMSP(username string) (map[string][]byte, error) {
	mspData := make(map[string][]byte)

	// Get IssuerRevocationPublicKey
	if content, err := v.getVaultContent(username, "IssuerRevocationPublicKey"); err == nil {
		mspData["IssuerRevocationPublicKey"] = content
	}

	// Get signcerts with error propagation
	if err := v.getDirectoryContent(username, "signcerts", mspData); err != nil {
		return nil, fmt.Errorf("signcerts error: %v", err)
	}

	// Get keystore with error propagation
	if err := v.getDirectoryContent(username, "keystore", mspData); err != nil {
		return nil, fmt.Errorf("keystore error: %v", err)
	}

	// Get cacerts with error propagation
	if err := v.getDirectoryContent(username, "cacerts", mspData); err != nil {
		return nil, fmt.Errorf("cacerts error: %v", err)
	}

	if len(mspData) == 0 {
		return nil, errors.New("no MSP data found")
	}

	return mspData, nil
}

func (v *VaultStore) getDirectoryContent(username, dir string, mspData map[string][]byte) error {
	// List directory contents using metadata path
	listPath := fmt.Sprintf("%s/metadata/users/%s/%s", v.kvPath, username, dir)
	fmt.Println(listPath)
	secret, err := v.client.Logical().List(listPath)
	if err != nil {
		return fmt.Errorf("failed to list directory: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return fmt.Errorf("directory not found")
	}

	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid directory structure")
	}

	for _, key := range keys {
		fileName := key.(string)
		// Skip directory markers
		if strings.HasSuffix(fileName, "/") {
			continue
		}

		content, err := v.getVaultContent(username, fmt.Sprintf("%s/%s", dir, fileName))
		if err != nil {
			return fmt.Errorf("failed to get %s: %w", fileName, err)
		}
		mspData[fmt.Sprintf("%s/%s", dir, fileName)] = content
	}
	return nil
}

func (v *VaultStore) getVaultContent(username, path string) ([]byte, error) {
	vaultPath := fmt.Sprintf("%s/data/users/%s/%s", v.kvPath, username, path)
	secret, err := v.client.Logical().Read(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("vault read failed: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("secret not found")
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format")
	}

	content, ok := data["content"].(string)
	if !ok {
		return nil, fmt.Errorf("content missing")
	}

	return []byte(content), nil
}
