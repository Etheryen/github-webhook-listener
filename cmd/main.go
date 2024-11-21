package main

import (
	"log"
	"net/http"
	"os"

	"github.com/etheryen/github-webhook-listener/internal/env"
	"github.com/etheryen/github-webhook-listener/pkg/webhook"
)

func main() {
	env.LoadEnv()

	http.HandleFunc("POST /webhook", webhook.Handler)

	listenPort := ":" + os.Getenv("PORT")

	log.Printf("Listening on %s\n", listenPort)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}
