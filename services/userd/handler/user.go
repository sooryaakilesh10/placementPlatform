package handler

import (
	"backend/services/userd/presenter"
	"backend/services/userd/usecase/user"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func getUserHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

func createUser(service user.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req presenter.UserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("unable to decode request body, err=%v", err), http.StatusBadRequest)
			return
		}

		newUser, err := service.CreateUser(req.UserName, req.Email, req.Pass, req.Role)
		if err == nil && newUser == nil {
			http.Error(w, fmt.Sprintf("unable to create user entity, err=%v", err), http.StatusBadRequest)
			return
		} else if err != nil {
			http.Error(w, fmt.Sprintf("unable to create user entity, err=%v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(presenter.UserResponse{
			UserID:   newUser.UserID,
			UserName: newUser.UserName,
			Email:    newUser.Email,
			Role:     newUser.Role,
		}); err != nil {
			http.Error(w, fmt.Sprintf("unable to encode to JSON, err=%v", err), http.StatusInternalServerError)
			return
		}
	}
}

func getUserByID(service user.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract ID from path /v1/user/id/{id}
		path := strings.TrimPrefix(r.URL.Path, "/v1/user/id/")
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
			http.Error(w, fmt.Sprintf("invalid user ID format: %v", err), http.StatusBadRequest)
			return
		}

		user, err := service.GetUserByID(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to get user by id, err=%v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(presenter.UserResponse{
			UserID:   user.UserID,
			UserName: user.UserName,
			Email:    user.Email,
			Role:     user.Role,
		}); err != nil {
			http.Error(w, fmt.Sprintf("unable to encode to JSON, err=%v", err), http.StatusInternalServerError)
			return
		}
	}
}

func getUserByEmail(service user.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract email from path /v1/user/email/user@example.com
		path := strings.TrimPrefix(r.URL.Path, "/v1/user/email/")
		if path == "" {
			http.Error(w, "email is required", http.StatusBadRequest)
			return
		}

		// Remove any trailing slashes and get the email
		email := strings.TrimSuffix(path, "/")
		if email == "" {
			http.Error(w, "email is required", http.StatusBadRequest)
			return
		}

		user, err := service.GetUserByEmail(email)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to get user by email, err=%v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(presenter.UserResponse{
			UserID:   user.UserID,
			UserName: user.UserName,
			Email:    user.Email,
			Role:     user.Role,
		}); err != nil {
			http.Error(w, fmt.Sprintf("unable to encode to JSON, err=%v", err), http.StatusInternalServerError)
			return
		}
	}
}

func login(service user.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req presenter.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("unable to decode request body, err=%v", err), http.StatusBadRequest)
			return
		}

		user, err := service.Login(req.Email, req.Pass)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to login user, err=%v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(presenter.LoginResponse{JWTToken: user}); err != nil {
			http.Error(w, fmt.Sprintf("unable to encode to JSON, err=%v", err), http.StatusInternalServerError)
			return
		}
	}
}

// Register User Routes
func RegisterUserHandlers(service user.Usecase) {
	http.HandleFunc("/v1/user/health", getUserHealth)           // GET
	http.HandleFunc("/v1/user", createUser(service))            // POST
	http.HandleFunc("/v1/user/id/", getUserByID(service))       //v1/user/id/{id} // GET
	http.HandleFunc("/v1/user/email/", getUserByEmail(service)) ////v1/user/id/{email} // GET
	http.HandleFunc("/v1/login", login(service))
}
