package main

import (
	"ZAPS/internal/handlers"
	"ZAPS/internal/queue"
	"ZAPS/internal/services"
	"ZAPS/internal/webhook"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("➡️ %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("⬅️ %s %s (%v)", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {

	q := queue.NewQueue(100)
	q.Start(5)

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
		q,
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

	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	mux.HandleFunc("/send-media", handlers.SendMedia(
		mediaService,
		messageService,
		conversationService,
	))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      loggingMiddleware(enableCORS(mux)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("🚀 Servidor rodando em http://localhost:8080")

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("🛑 Desligando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("Erro ao desligar servidor:", err)
	} else {
		log.Println("✅ Servidor finalizado com sucesso")
	}
}
