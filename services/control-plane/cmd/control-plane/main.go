package main

import (
	"log"
	"net/http"
	"os"

	"github.com/beparykamrul-dev/FTN-AI/services/control-plane/internal/api"
)

func main() {
	addr := os.Getenv("FTN_CONTROL_PLANE_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	server := &http.Server{
		Addr:    addr,
		Handler: api.Router(),
	}

	log.Printf("FTN control-plane listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
