package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/koderkt/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")
	connString := os.Getenv("CONN_STR")
	if connString == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("Could not connect to db: ", err)
	}

	dbQueries := database.New(db)

	apiConfig := &apiConfig{
		DB: dbQueries,
	}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", handlerReadiness)
	mux.HandleFunc("GET /v1/err", handlerErr)
	mux.HandleFunc("POST /v1/users", apiConfig.handlerCreateUser)
	mux.HandleFunc("GET /v1/users", apiConfig.middlewareAuth(apiConfig.handlerGetUser))
	mux.HandleFunc("POST /v1/feeds", apiConfig.middlewareAuth(apiConfig.handlerCreateFeed))
	mux.HandleFunc("GET /v1/feeds", apiConfig.handlerGetAllFeeds)
	// TODO: TEST
	mux.HandleFunc("POST /v1/feed_follows", apiConfig.middlewareAuth(apiConfig.handlerFeedFollowCreate))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiConfig.middlewareAuth(apiConfig.handlerDeleteFollow))
	mux.HandleFunc("GET /v1/feed_follows", apiConfig.middlewareAuth(apiConfig.handlerGetAllFeedFollowsForAUser))
	mux.HandleFunc("GET /v1/posts", apiConfig.middlewareAuth(apiConfig.handlerPostsGet))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	go startScraper(apiConfig.DB, time.Minute)
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
