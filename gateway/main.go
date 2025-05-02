package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"log"
	"net/http"
)

type Identity struct {
	Id               string `json:"id"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	Dob              string `json:"dob"`
	PresentAddress   string `json:"presentAddress"`
	PermanentAddress string `json:"permanentAddress"`
	Gender           string `json:"gender"`
	NationalID       string `json:"nationalID"`
	Owner            string `json:"owner"`
}

type FabricClient struct {
	gateway  *client.Gateway
	contract *client.Contract
}

func main() {
	// Initialize Fabric connection
	fabricClient, err := InitFabricConnection()
	if err != nil {
		log.Fatalf("Failed to initialize Fabric connection: %v", err)
	}
	defer func() {
		err = fabricClient.gateway.Close()
		log.Fatalf("Failed to close gateway connection: %v", err)
	}()

	// Create router
	r := chi.NewRouter()

	// Set up middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Set up routes
	r.Post("/create", fabricClient.createIdentity)
	r.Post("/update", fabricClient.updateIdentity)
	r.Post("/delete", fabricClient.deleteIdentity)
	r.Get("/get/{id}", fabricClient.getIdentity)

	// Start server
	port := envOrDefault("PORT", "8080")
	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func (fc *FabricClient) createIdentity(w http.ResponseWriter, r *http.Request) {
	var idnty Identity
	if err := json.NewDecoder(r.Body).Decode(&idnty); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if isEmptyField(idnty.Id) {
		http.Error(w, "identity id is not provided", http.StatusBadRequest)
		return
	}

	if isEmptyField(idnty.FirstName) {
		http.Error(w, "identity firstName is not provided", http.StatusBadRequest)
		return
	}

	if isEmptyField(idnty.Phone) {
		http.Error(w, "identity phone is not provided", http.StatusBadRequest)
	}

	if isEmptyField(idnty.NationalID) {
		http.Error(w, "identity national id is not provided", http.StatusBadRequest)
	}

	assetJSON, err := json.Marshal(idnty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := fc.contract.SubmitTransaction("CreateAsset", string(assetJSON)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"AssetId": idnty.Id})
}

func (fc *FabricClient) updateIdentity(w http.ResponseWriter, r *http.Request) {
	var idnty Identity
	if err := json.NewDecoder(r.Body).Decode(&idnty); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if isEmptyField(idnty.Id) {
		http.Error(w, "identity id is not provided", http.StatusBadRequest)
		return
	}

	assetJSON, err := json.Marshal(idnty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := fc.contract.SubmitTransaction("UpdateAsset", idnty.Id, string(assetJSON)); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Update success",
	})
}

func (fc *FabricClient) deleteIdentity(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if isEmptyField(req.ID) {
		http.Error(w, "identity id is not provided", http.StatusBadRequest)
		return
	}

	if _, err := fc.contract.SubmitTransaction("DeleteAsset", req.ID); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong: " + err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Delete success",
	})
}

func (fc *FabricClient) getIdentity(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if isEmptyField(id) {
		http.Error(w, "identity id is not provided", http.StatusBadRequest)
		return
	}

	result, err := fc.contract.EvaluateTransaction("ReadAsset", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var idnty Identity
	if err := json.Unmarshal(result, &idnty); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, idnty)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func isEmptyField(f string) bool {
	if f == "" {
		return true
	}
	return false
}
