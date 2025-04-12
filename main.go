package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	// _ "github.com/lib/pq"

	"backend/pkg/common"
	storageHandler "backend/services/storaged/handler"
	storageRepo "backend/services/storaged/repository"
	storageUsecase "backend/services/storaged/usecase/storage"
	userHandler "backend/services/userd/handler"
	userRepo "backend/services/userd/repository"
	userUsecase "backend/services/userd/usecase/user"
)

const PORT = "8080"

func main() {
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		getEnv("DB_USER", "mysqluser"),
		getEnv("DB_PASSWORD", "mysqlpassword"),
		getEnv("DB_HOST", "mysql"),
		getEnv("DB_NAME", "portal"))

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	for i := 0; i < 30 && db.Ping() != nil; i++ {
		if i == 29 {
			log.Fatalf("Database not available after 30 attempts")
		}
		log.Printf("Database not ready, waiting... (attempt %d/30)", i+1)
		time.Sleep(2 * time.Second)
	}
	log.Println("Database connection successful")

	// Initialize system settings
	common.InitSystemSettings(db)

	// Initialize repositories
	userRepository := userRepo.NewRepository(db)
	storageRepository := storageRepo.NewRepository(db)

	// Initialize services
	userService := userUsecase.NewService(userRepository)
	storageService := storageUsecase.NewService(storageRepository)

	// Register handlers
	userHandler.RegisterUserHandlers(userService)
	storageHandler.RegisterCompanyHandlers(storageService)

	// Set up CORS middleware
	http.HandleFunc("/", corsMiddleware(http.DefaultServeMux))

	port := getEnv("PORT", PORT)
	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

// corsMiddleware handles CORS for all routes
func corsMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Process the request
		next.ServeHTTP(w, r)
	}
}
