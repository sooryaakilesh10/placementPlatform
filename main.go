package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	// _ "github.com/lib/pq"

	userHandler "backend/services/userd/handler"
	"backend/services/userd/repository"
	"backend/services/userd/usecase/user"
)

const PORT = "8080"

func generateSecret(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

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

	jwtSecret, err := generateSecret(32)
	if err != nil {
		log.Fatalf("Error generating JWT secret: %v", err)
	}
	userHandler.RegisterUserHandlers(user.NewService(repository.NewRepository(db), jwtSecret))

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
