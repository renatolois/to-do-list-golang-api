package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/renatolois/to-do-list-golang-api/internal/task"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "tasks.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	repo, err := task.NewSQLiteRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	srv := task.NewService(repo)
	h := task.NewHandler(srv)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", h.HandleList)
	mux.HandleFunc("POST /tasks", h.HandleCreate)
	mux.HandleFunc("GET /tasks/{id}", h.HandleGetByID)
	mux.HandleFunc("DELETE /tasks/{id}", h.HandleDelete)

	log.Println("Server is running on http://localhost:8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
