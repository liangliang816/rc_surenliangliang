package main

import (
	"log"

	"github.com/gin-gonic/gin/controller"
	"github.com/gin-gonic/gin/worker"
)

func main() {
	// Start background worker
	worker.StartWorker()
	log.Println("Background worker started")

	// Setup router with API endpoints
	r := controller.SetupRouter()

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}