package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/etheryen/github-webhook-listener/internal/env"
	"github.com/etheryen/github-webhook-listener/pkg/webhook"
)

func main() {
	githubSecret, port, configPath := env.GetVars()

	config := webhook.ParseConfig(configPath)

	fmt.Println()
	config.Print()
	fmt.Println()

	http.HandleFunc("POST /webhook", webhook.Handler(config, githubSecret))

	listenPort := ":" + port

	log.Printf("Listening on %s\n", listenPort)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}
