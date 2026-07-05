package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"short-track/internal/handler"
	"short-track/internal/repository"
	"short-track/internal/shortener"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	// 1. Load configuration from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "pass")
	dbName := getEnv("DB_NAME", "shortener")

	redisHost := getEnv("REDIS_HOST", "localhost:6379")

	// 2. Initialize PostgreSQL Connection with retry logic
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	log.Printf("Connecting to Postgres on %s:%s...", dbHost, dbPort)
	
	var db *sql.DB
	var err error
	for i := 0; i < 15; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Failed to connect to database (attempt %d/15): %v. Retrying in 2 seconds...", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	// 3. Initialize Redis Client with retry logic
	log.Printf("Connecting to Redis on %s...", redisHost)
	rdb := redis.NewClient(&redis.Options{
		Addr: redisHost,
	})
	
	var redisErr error
	for i := 0; i < 15; i++ {
		_, redisErr = rdb.Ping(rdb.Context()).Result()
		if redisErr == nil {
			break
		}
		log.Printf("Failed to connect to Redis (attempt %d/15): %v. Retrying in 2 seconds...", i+1, redisErr)
		time.Sleep(2 * time.Second)
	}
	if redisErr != nil {
		log.Fatalf("Could not connect to Redis: %v", redisErr)
	}
	defer rdb.Close()

	// 4. Dependency Injection
	repo := repository.NewRepository(db, rdb)
	svc := shortener.NewService(repo)
	h := handler.NewHandler(svc)

	// 5. Routing
	r := gin.Default()
	r.POST("/shorten", h.ShortenURL)
	r.GET("/:code", h.Redirect)

	log.Println("Starting API server on :8080...")
	r.Run(":8080")
}