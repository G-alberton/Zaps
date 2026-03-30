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
	mediaService := services.NewMediaService(
		"EAARpsmqfXFUBRLCsLK8L2e3rI7vY24aQn2nBDlrgylkWPTomlHWQyMOkBGFupJbEwZBODQ97j5PqMymD4SMURLAVZAYGkVCIoXTVwgzLrZAXZAW9Ft8LJ116OYfatgYKCvwXkMuy0mKbJ8fLyMNy2s91Ulj6SsrAqY4DvIYmlY87ylWhttoubxJeOPqQfvnv01NvpMJpt5zKE8JaTR8rtDUbnGjg4VBIfM6auTW6JKuWCXV96kFR3srhZA0z09jM1IiElhoDvEoRvLlZBZBM3cd",
		"991625617368927",
	)
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
	mux.HandleFunc("/send-message", handlers.SendMessage(
		mediaService,
		messageService,
		conversationService,
	))

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
