package main

import (
	"log"
	"net/http"

	"ZAPS/internal/handlers"
	"ZAPS/internal/services"
	"ZAPS/internal/webhook"
)

func main() {
	contactService := services.NewContactService(nil)
	messageService := services.NewMessageService(nil)
	mediaService := services.NewMediaService("Token")
	conversationService := services.NewConversationService()

	mux := http.NewServeMux()

	// webhook
	mux.HandleFunc("/webhook", webhook.HandleWebhook(
		contactService,
		messageService,
		mediaService,
		conversationService,
	))

	// conversations
	mux.HandleFunc("/conversations", handlers.GetConversations(conversationService))

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
