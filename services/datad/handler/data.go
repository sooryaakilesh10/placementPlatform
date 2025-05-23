package handler

import (
	"backend/services/datad/presenter"
	dataRepository "backend/services/datad/repository"
	"backend/services/datad/usecase/data"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func getDataHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func createCompany(service data.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req presenter.CreateCompanyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Unable to decode request body, err=%v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		compnayID, err := service.CreateCompany(
			req.JWT,
			req.CompanyName,
			req.CompanyAddress,
			req.Drive,
			req.TypeOfDrive,
			req.FollowUp,
			req.Remarks,
			req.ContactDetails,
			req.HRDetails,
			req.IsContacted)
		if err != nil {
			log.Printf("Unable to create company, err=%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(presenter.CreateCompanyResponse{CompanyID: compnayID}); err != nil {
			log.Printf("Unable to encode response, err=%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func getCompany(service data.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract ID from path /v1/data/id/{id}
		path := strings.TrimPrefix(r.URL.Path, "/v1/data/id/")
		if path == "" {
			http.Error(w, "company ID is required in the path", http.StatusBadRequest)
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
			http.Error(w, fmt.Sprintf("invalid company ID format: %v", err), http.StatusBadRequest)
			return
		}

		// TODO: Implement proper JWT extraction and validation if needed for GET
		jwtString := "" // Placeholder: Pass empty JWT for now

		company, err := service.GetCompany(jwtString, id) // Use the 'id' variable from the path
		if err != nil {
			if errors.Is(err, dataRepository.ErrNotFound) {
				http.Error(w, "Company not found", http.StatusNotFound)
			} else {
				log.Printf("Unable to get company, err=%v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(presenter.GetCompanyResponse{
			CompanyID:      company.CompanyID,
			CompanyName:    company.CompanyName,
			CompanyAddress: company.CompanyAddress,
			Drive:          company.Drive,
			TypeOfDrive:    company.TypeOfDrive,
			FollowUp:       company.FollowUp,
			IsContacted:    company.IsContacted,
			Remarks:        company.Remarks,
			ContactDetails: company.ContactDetails,
			HRDetails:      company.HRDetails,
		}); err != nil {
			log.Printf("Unable to encode response, err=%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func getCompanyByName(service data.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract ID from path /v1/data/id/{id}
		path := strings.TrimPrefix(r.URL.Path, "/v1/data/name/")
		if path == "" {
			http.Error(w, "company ID is required in the path", http.StatusBadRequest)
			return
		}

		// Remove any trailing slashes and get the ID
		name := strings.TrimSuffix(path, "/")
		if name == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		// TODO: Implement proper JWT extraction and validation if needed for GET
		jwtString := "" // Placeholder: Pass empty JWT for now

		company, err := service.GetCompany(jwtString, name) // Use the 'id' variable from the path
		if err != nil {
			if errors.Is(err, dataRepository.ErrNotFound) {
				http.Error(w, "Company not found", http.StatusNotFound)
			} else {
				log.Printf("Unable to get company, err=%v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(presenter.GetCompanyResponse{
			CompanyID:      company.CompanyID,
			CompanyName:    company.CompanyName,
			CompanyAddress: company.CompanyAddress,
			Drive:          company.Drive,
			TypeOfDrive:    company.TypeOfDrive,
			FollowUp:       company.FollowUp,
			IsContacted:    company.IsContacted,
			Remarks:        company.Remarks,
			ContactDetails: company.ContactDetails,
			HRDetails:      company.HRDetails,
		}); err != nil {
			log.Printf("Unable to encode response, err=%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func updateCompanyByID(service data.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract ID from path /v1/data/id/{id}
		path := strings.TrimPrefix(r.URL.Path, "/v1/data/name/")
		if path == "" {
			http.Error(w, "company ID is required in the path", http.StatusBadRequest)
			return
		}

		// Remove any trailing slashes and get the ID
		id := strings.TrimSuffix(path, "/")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		var req presenter.CreateCompanyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Unable to decode request body, err=%v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		company, err := service.UpdateCompany(
			req.JWT,
			req.CompanyID,
			req.CompanyName,
			req.CompanyAddress,
			req.Drive,
			req.TypeOfDrive,
			req.FollowUp,
			req.Remarks,
			req.ContactDetails,
			req.HRDetails,
			req.IsContacted)
		if err != nil {
			log.Printf("Unable to create company, err=%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(presenter.GetCompanyResponse{
			CompanyID:      company.CompanyID,
			CompanyName:    company.CompanyName,
			CompanyAddress: company.CompanyAddress,
			Drive:          company.Drive,
			TypeOfDrive:    company.TypeOfDrive,
			FollowUp:       company.FollowUp,
			IsContacted:    company.IsContacted,
			Remarks:        company.Remarks,
			ContactDetails: company.ContactDetails,
			HRDetails:      company.HRDetails,
		}); err != nil {
			log.Printf("Unable to encode response, err=%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func getAwaitingApproval(service data.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract ID from path /v1/data/id/{id}
		path := strings.TrimPrefix(r.URL.Path, "/v1/data/name/")
		if path == "" {
			http.Error(w, "company ID is required in the path", http.StatusBadRequest)
			return
		}

		// Remove any trailing slashes and get the ID
		id := strings.TrimSuffix(path, "/")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		companies, err := service.GetAwaitingApproval()
		if err != nil {
			log.Printf("Unable to create company, err=%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(companies); err != nil {
			log.Printf("Unable to encode response, err=%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func setAwaitingApproval(service data.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req presenter.SetAwaitingApprovalRequest
		// Extract ID from path /v1/data/approve/id/{id}
		path := strings.TrimPrefix(r.URL.Path, "/v1/data/approve/id/")
		if path == "" {
			http.Error(w, "company ID is required in the path", http.StatusBadRequest)
			return
		}

		// Remove any trailing slashes and get the ID
		id := strings.TrimSuffix(path, "/")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Unable to decode request body, err=%v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err := service.SetAwaitingApproval(req.JWT, id, req.IsApproved)
		if err != nil {
			if errors.Is(err, dataRepository.ErrNotFound) {
				http.Error(w, "Company not found", http.StatusNotFound)
			} else {
				log.Printf("Unable to create company, err=%v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Service call is successful, return 200 OK with empty body
		w.WriteHeader(http.StatusOK)
		// Optional: Write an empty JSON object {} or just leave the body empty
		// w.Write([]byte("{}"))
	}
}

// Register Data Routes
func RegisterDataHandlers(service data.Usecase) {
	http.HandleFunc("/v1/data/health", getDataHealth)   // GET
	http.HandleFunc("/v1/data", createCompany(service)) // POST
	http.HandleFunc("/v1/data/id/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCompany(service)(w, r) // GET
		case http.MethodPut:
			updateCompanyByID(service)(w, r) // PUT
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/v1/data/name/", getCompanyByName(service))          // POST
	http.HandleFunc("/v1/data/approve", getAwaitingApproval(service))     // GET
	http.HandleFunc("/v1/data/approve/id/", setAwaitingApproval(service)) // GET
}
