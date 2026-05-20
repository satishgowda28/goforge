package main

import (
	"context"
	"fmt"
	"log"

	"github.com/satishgowda28/goforge/internal/agent"
	"github.com/satishgowda28/goforge/internal/llm/anthropic"
	"github.com/satishgowda28/goforge/internal/tools"
	"github.com/satishgowda28/goforge/pkg/config"
)

func main() {
	_, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	// llm provider
	client := anthropic.NewClient()

	//tool registry
	t := tools.NewRegistry()
	t.Register(&tools.HttpFetch{})
	t.Register(&tools.FileWrite{})

	//agent registry
	aRegistry, err := agent.NewAgentRegistry("./agents")
	if err != nil {
		log.Fatalf("Error in Registering agents: %v", err)
	}
	researcherConfig, err := aRegistry.Get("researcher")
	if err != nil {
		log.Fatalf("%v", err)
	}

	// memeory
	var memory agent.ShortTermMemory

	runner := agent.NewRunner(client, t, &memory, researcherConfig)
	result, err := runner.Run(context.Background(), "Whats happeining new in AI")
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Print(result)
	// r := api.InitServer()
	// port := fmt.Sprintf(":%d", cfg.Server.Port)

	// if err := http.ListenAndServe(port, r); err != nil {
	// 	log.Fatalf("Error in Starting a server: %v", err)
	// }
}
