package main

import (
	"log"
	"net/http"

	"ZAPS/internal/database"
	"ZAPS/internal/repository"
	"ZAPS/internal/services"
	"ZAPS/internal/webhook"
)

func main() {
	db := database.Connect()

	contactRepo := repository.NewContactRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	contactService := services.NewContactService(contactRepo)
	messageService := services.NewMessageService(messageRepo)
	mediaService := services.NewMediaService("Token")

	mux := http.NewServeMux()
	// passa serviço para o webhook
	mux.HandleFunc("/webhook", webhook.HandleWebhook(contactService, messageService, mediaService))

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
