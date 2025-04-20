package handler

import (
	"backend/services/datad/presenter"
	"backend/services/datad/usecase/data"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func getDataHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func createData(service data.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req presenter.DataRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "unable to decode request body", http.StatusBadRequest)
			return
		}

		// If JWT is not in request body, try to get it from Authorization header
		if req.JWT == "" {
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				req.JWT = authHeader[7:]
			}
		}

		// Validate company_data is not nil
		if req.CompanyData == nil {
			http.Error(w, "company_data cannot be empty", http.StatusBadRequest)
			return
		}

		newData, err := service.CreateData(req.JWT, req.CompanyData)
		if err != nil {
			statusCode := http.StatusInternalServerError

			// Determine appropriate status code based on error message
			errMsg := err.Error()
			if errMsg == "permission denied: insufficient role" {
				statusCode = http.StatusForbidden
			} else if errMsg == "invalid token" || (len(errMsg) > 13 && errMsg[:13] == "invalid token:") {
				statusCode = http.StatusBadRequest
			}

			http.Error(w, fmt.Sprintf("unable to create data entity: %v", err), statusCode)
			return
		}

		if newData == nil {
			http.Error(w, "unexpected server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(struct {
			ID string `json:"data_id"`
		}{ID: newData.DataID}); err != nil {
			http.Error(w, "unable to encode to JSON", http.StatusInternalServerError)
			return
		}
	}
}

func getDataByID(service data.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract ID from path /v1/data/id/{id}
		path := strings.TrimPrefix(r.URL.Path, "/v1/data/id/")
		if path == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		// Remove any trailing slashes and get the ID
		id := strings.TrimSuffix(path, "/")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		// Validate if the ID is a valid UUID
		_, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid data ID format: %v", err), http.StatusBadRequest)
			return
		}

		// Get data by ID
		data, err := service.GetDataByID(id)
		if err != nil {
			http.Error(w, "unable to get data by ID", http.StatusInternalServerError)
			return
		}

		if data == nil {
			http.Error(w, "data not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(
			presenter.DataResponse{DataID: data.DataID, CompanyData: data.CompanyData},
		); err != nil {
			http.Error(w, "unable to encode to JSON", http.StatusInternalServerError)
			return
		}
	}
}

// func getDataByName(service data.Usecase) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodGet {
// 			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		vars := mux.Vars(r)
// 		name := vars["name"]

// 		log.Printf("Received data retrieval request for name: %s", name)

// 		// Validate name
// 		if name == "" {
// 			http.Error(w, "name is empty", http.StatusBadRequest)
// 			return
// 		}

// 		// Get data by name
// 		data, err := service.GetDataByName(name)
// 		if err != nil {
// 			http.Error(w, "unable to get data by name", http.StatusInternalServerError)
// 			return
// 		}

// 		if data == nil {
// 			http.Error(w, "data not found", http.StatusNotFound)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		if err := json.NewEncoder(w).Encode(presenter.DataResponse{DataID: data.DataID}); err != nil {
// 			http.Error(w, fmt.Sprintf("unable to encode to JSON", err), http.StatusInternalServerError)
// 			return
// 		}
// 	}
// }

// func listData(service user.Usecase) http.HandlerFunc {
// }

// func login(service user.Usecase) http.HandlerFunc {
// }

// Register Data Routes
func RegisterDataHandlers(service data.Usecase) {
	http.HandleFunc("/v1/data/health", getDataHealth)     // GET
	http.HandleFunc("/v1/data", createData(service))      // POST
	http.HandleFunc("/v1/data/id/", getDataByID(service)) // GET /v1/data/id/{id}
	// http.HandleFunc("/v1/data/name/", getDataByName(service)) // GET /v1/data/name/{name}
}
