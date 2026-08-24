package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"learn-golang/internal/api"
	"learn-golang/internal/content"
	"learn-golang/internal/progress"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	contentRoot := filepath.Join(root, "content", "go-zero-to-hero")
	webRoot := filepath.Join(root, "web")
	dbPath := filepath.Join(root, "data", "progress.db")

	contentStore, err := content.NewStore(contentRoot)
	if err != nil {
		log.Fatalf("load content: %v", err)
	}

	progressStore, err := progress.NewStore(dbPath)
	if err != nil {
		log.Fatalf("init progress: %v", err)
	}
	defer progressStore.Close()

	addr := envOr("ADDR", ":8080")
	server := api.NewServer(contentStore, progressStore, webRoot)

	log.Printf("Learn Go server running at http://localhost%s", addr)
	if err := http.ListenAndServe(addr, server.Router()); err != nil {
		log.Fatal(err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
