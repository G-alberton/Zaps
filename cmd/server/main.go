package main

import (
	"log"
	"net/http"

	"ZAPS/internal/handlers"
	"ZAPS/internal/services"
	"ZAPS/internal/webhook"
)

func main() {

	mediaService := services.NewMediaService()
	conversationService := services.NewConversationService()
	messageService := services.NewMessageService(nil) // sem banco por enquanto
	contactService := services.NewContactService(nil) // sem banco por enquanto

	mux := http.NewServeMux()

	mux.HandleFunc("/webhook", webhook.HandleWebhook(
		contactService,
		messageService,
		mediaService,
		conversationService,
	))

	mux.HandleFunc("/send-message", handlers.SendMessage(
		mediaService,
		messageService,
		conversationService,
	))

	mux.HandleFunc("/messages", handlers.GetMessages(messageService))

	mux.HandleFunc("/conversations", handlers.GetConversations(conversationService, messageService))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("🚀 Servidor rodando em http://localhost:8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
