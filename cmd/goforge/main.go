package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/satishgowda28/goforge/internal/api"
	"github.com/satishgowda28/goforge/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	r := api.InitServer()
	port := fmt.Sprintf(":%d", cfg.Server.Port)

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Error in Starting a server: %v", err)
	}
}
