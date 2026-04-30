package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"ZAPS/internal/auth"
	"ZAPS/internal/database"
	"ZAPS/internal/handlers"
	"ZAPS/internal/middleware"
	"ZAPS/internal/queue"
	"ZAPS/internal/repository"
	"ZAPS/internal/services"
	"ZAPS/internal/webhook"
	"ZAPS/internal/websocket"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Message-ID")
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

	if err := godotenv.Load(); err != nil {
		log.Println(".env não encontrado")
	}

	db := database.Connect()

	userRepo := &repository.PostgresUserRepository{DB: db}
	messageRepo := repository.NewMessageRepository(db)
	conversationRepo := repository.NewConversationRepository(db)

	jwtService := &auth.JWTService{
		Secret: []byte("super-secret"),
		Expire: 24 * time.Hour,
	}

	authMiddleware := middleware.AuthMiddleware(jwtService)

	authService := &services.AuthService{
		Repo: userRepo,
		JWT:  jwtService,
	}

	authHandler := &handlers.AuthHandler{
		Service: authService,
	}

	hub := websocket.NewHub()
	go hub.Run()

	q := queue.NewPriorityQueue(100)
	q.Start(5)

	mediaService := services.NewMediaService()
	conversationService := services.NewConversationService(conversationRepo)

	messageService := services.NewMessageService(
		messageRepo,
		conversationRepo,
		hub,
	)

	sendService := services.NewSendService(
		mediaService,
		messageService,
		hub,
	)

	contactService := services.NewContactService(nil)

	conversationHandler := handlers.NewConversationHandler(
		conversationService,
		messageService,
		contactService,
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		tokenStr := r.URL.Query().Get("token")
		if tokenStr == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		claims, err := jwtService.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		conversationID := r.URL.Query().Get("conversation_id")
		if conversationID == "" {
			http.Error(w, "conversation_id required", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "conversation_id", conversationID)

		log.Println("WS conectado | user_id:", claims.UserID, "| conversation:", conversationID)

		websocket.ServerWS(hub, w, r.WithContext(ctx))
	})

	mux.HandleFunc("/webhook", webhook.HandleWebhook(
		contactService,
		messageService,
		mediaService,
		conversationService,
		q,
		hub,
	))

	mux.Handle("/", http.FileServer(http.Dir("../../web")))

	mux.Handle("/messages",
		authMiddleware(http.HandlerFunc(handlers.GetMessages(messageService))),
	)

	mux.Handle("/messages/paginated",
		authMiddleware(http.HandlerFunc(handlers.ListMessagesPaginated(messageService))),
	)

	mux.Handle("/mark-as-read",
		authMiddleware(http.HandlerFunc(handlers.MarkAsRead(messageService))),
	)

	mux.Handle(
		"/send-message",
		authMiddleware(
			handlers.SendMessage(
				mediaService,
				messageService,
				conversationService,
				hub,
			),
		),
	)

	mux.Handle(
		"/send-media",
		authMiddleware(
			handlers.SendMedia(
				messageService,
				conversationService,
				sendService,
				hub,
			),
		),
	)

	mux.Handle("/conversations",
		authMiddleware(http.HandlerFunc(conversationHandler.GetConversations)),
	)

	mux.Handle("/conversation",
		authMiddleware(http.HandlerFunc(conversationHandler.GetOrCreate)),
	)

	mux.Handle("/uploads/",
		authMiddleware(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads")))),
	)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      loggingMiddleware(enableCORS(mux)),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
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
