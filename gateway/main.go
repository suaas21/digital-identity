package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"google.golang.org/grpc"
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

func main() {
	// new grpc connection to fabric peer
	grpcConn, err := newGrpcConnection()
	if err != nil {
		log.Fatalf("Failed to initialize gRPC connection: %v", err)
	}
	defer func() {
		err = grpcConn.Close()
		log.Fatalf("Failed to close GRPC connection: %v", err)
	}()

	// Create router
	r := chi.NewRouter()

	// Set up middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-User-Cert", "X-User-Key", "X-User-MSPID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Set up routes
	r.Post("/create", createIdentityHandler(grpcConn))
	r.Post("/update", updateIdentityHandler(grpcConn))
	r.Post("/delete", deleteIdentityHandler(grpcConn))
	r.Get("/get/{id}", getIdentityHandler(grpcConn))

	// Start server
	port := envOrDefault("PORT", "8080")
	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func createIdentityHandler(grpcConn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract identity headers
		certPEM := r.Header.Get("X-User-Cert")
		keyPEM := r.Header.Get("X-User-Key")
		mspID := r.Header.Get("X-User-MSPID")

		if certPEM == "" || keyPEM == "" || mspID == "" {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "Missing required identity headers (X-User-Cert, X-User-Key, X-User-MSPID)",
			})
			return
		}

		// Create gateway connection for this identity
		gw, contract, err := newGatewayFromIdentity(grpcConn, certPEM, keyPEM, mspID)
		if err != nil {
			respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"status":  http.StatusUnauthorized,
				"message": "Invalid identity credentials: " + err.Error(),
			})
			return
		}
		defer gw.Close()

		// Parse request body
		var idnty Identity
		if err := json.NewDecoder(r.Body).Decode(&idnty); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body: " + err.Error(),
			})
			return
		}

		// Validate required fields
		if isEmptyField(idnty.Id) {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "Identity ID is required",
			})
			return
		}

		requiredFields := map[string]string{
			"firstName":  idnty.FirstName,
			"phone":      idnty.Phone,
			"nationalID": idnty.NationalID,
		}

		for field, value := range requiredFields {
			if isEmptyField(value) {
				respondJSON(w, http.StatusBadRequest, map[string]interface{}{
					"status":  http.StatusBadRequest,
					"message": fmt.Sprintf("Field %s is required", field),
				})
				return
			}
		}

		// Submit transaction
		assetJSON, err := json.Marshal(idnty)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"status":  http.StatusInternalServerError,
				"message": "Error marshaling identity: " + err.Error(),
			})
			return
		}

		if _, err := contract.SubmitTransaction("CreateIdentity", string(assetJSON)); err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"status":  http.StatusInternalServerError,
				"message": "Chaincode error: " + err.Error(),
			})
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"status":  http.StatusOK,
			"message": "Identity created successfully",
			"assetId": idnty.Id,
		})
	}
}

func updateIdentityHandler(grpcConn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract identity headers
		certPEM := r.Header.Get("X-User-Cert")
		keyPEM := r.Header.Get("X-User-Key")
		mspID := r.Header.Get("X-User-MSPID")

		if certPEM == "" || keyPEM == "" || mspID == "" {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "Missing required identity headers",
			})
			return
		}

		// Create gateway connection
		gw, contract, err := newGatewayFromIdentity(grpcConn, certPEM, keyPEM, mspID)
		if err != nil {
			respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"status":  http.StatusUnauthorized,
				"message": "Invalid identity credentials: " + err.Error(),
			})
			return
		}
		defer gw.Close()

		// Parse request
		var idnty Identity
		if err = json.NewDecoder(r.Body).Decode(&idnty); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body:  " + err.Error(),
			})
			return
		}

		if isEmptyField(idnty.Id) {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "Identity ID is required",
			})
			return
		}

		// Submit transaction
		assetJSON, err := json.Marshal(idnty)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"status":  http.StatusInternalServerError,
				"message": "Error marshaling identity:  " + err.Error(),
			})
			return
		}

		if _, err = contract.SubmitTransaction("UpdateIdentity", idnty.Id, string(assetJSON)); err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"status":  http.StatusInternalServerError,
				"message": "Chaincode error: " + err.Error(),
			})
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"status":  http.StatusOK,
			"message": "Identity updated successfully",
			"assetId": idnty.Id,
		})
	}
}

func deleteIdentityHandler(grpcConn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract identity headers
		certPEM := r.Header.Get("X-User-Cert")
		keyPEM := r.Header.Get("X-User-Key")
		mspID := r.Header.Get("X-User-MSPID")

		if certPEM == "" || keyPEM == "" || mspID == "" {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "Missing required identity headers",
			})
			return
		}

		// Create gateway connection
		gw, contract, err := newGatewayFromIdentity(grpcConn, certPEM, keyPEM, mspID)
		if err != nil {
			respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"status":  http.StatusUnauthorized,
				"message": "Invalid identity credentials:  " + err.Error(),
			})
			return
		}
		defer gw.Close()

		// Parse request
		var request struct {
			ID string `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
			})
			return
		}

		if isEmptyField(request.ID) {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "Identity ID is required",
			})
			return
		}

		// Submit transaction
		if _, err = contract.SubmitTransaction("DeleteIdentity", request.ID); err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"status":  http.StatusInternalServerError,
				"message": "Chaincode error: " + err.Error(),
			})
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"status":  http.StatusOK,
			"message": "Identity deleted successfully",
			"assetId": request.ID,
		})
	}
}

func getIdentityHandler(grpcConn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract identity headers
		certPEM := r.Header.Get("X-User-Cert")
		keyPEM := r.Header.Get("X-User-Key")
		mspID := r.Header.Get("X-User-MSPID")

		if certPEM == "" || keyPEM == "" || mspID == "" {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "Missing required identity headers",
			})
			return
		}

		// Create gateway connection
		gw, contract, err := newGatewayFromIdentity(grpcConn, certPEM, keyPEM, mspID)
		if err != nil {
			respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"status":  http.StatusUnauthorized,
				"message": fmt.Sprintf("Invalid identity credentials: err: %v", err),
			})
			return
		}
		defer gw.Close()

		// Get ID from URL
		id := chi.URLParam(r, "id")
		if isEmptyField(id) {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "Identity ID is required",
			})
			return
		}

		// Evaluate transaction
		result, err := contract.EvaluateTransaction("ReadIdentity", id)
		if err != nil {
			respondJSON(w, http.StatusNotFound, map[string]interface{}{
				"status":  http.StatusNotFound,
				"message": "Identity not found: " + err.Error(),
			})
			return
		}

		var identity Identity
		if err = json.Unmarshal(result, &identity); err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"status":  http.StatusInternalServerError,
				"message": "Error parsing identity data:  " + err.Error(),
			})
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"status":   http.StatusOK,
			"identity": identity,
		})
	}
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
