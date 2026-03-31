package main

import (
	"log"
	"net/http"
	"time"

	"ZAPS/internal/handlers"
	"ZAPS/internal/services"
	"ZAPS/internal/webhook"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

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

	mux.HandleFunc("/conversations", handlers.GetConversations(
		conversationService,
		messageService,
		contactService,
	))

	mux.HandleFunc("/mark-as-read", handlers.MarkAsRead(messageService))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("🚀 Servidor rodando em http://localhost:8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
