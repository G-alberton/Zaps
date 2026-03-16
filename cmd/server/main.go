package main

import (
	"log"
	"net/http"

	"ZAPS/internal/webhook"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/webhook", webhook.HandleWebhook)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Servidor rodando na porta 8080")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
