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

	userHandler "backend/services/userd/handler"
	"backend/services/userd/repository"
	"backend/services/userd/usecase/user"
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

	userHandler.RegisterUserHandlers(user.NewService(repository.NewRepository(db)))

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
