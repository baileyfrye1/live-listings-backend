package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"server/database"
	"server/internal/api/handler"
	"server/internal/logger"
	"server/internal/repo"
	"server/internal/server"
	"server/internal/service"
	"server/internal/session"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	slog.Info("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown with error: ", "error", err)
	}

	slog.Info("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	logger.Init(logger.Config{
		LogLevel:   slog.LevelDebug,
		JSONFormat: false,
	})

	dbService := database.New()

	// Initialize Redis client and session
	client := session.GetClient()
	session := session.NewSession(client)

	// Setup repositories
	userRepo := repo.NewUserRepository(dbService.DB())
	listingRepo := repo.NewListingRepository(dbService.DB())

	// Setup services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, session)
	listingService := service.NewListingService(listingRepo)

	// Setup handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	listingHandler := handler.NewListingHandler(listingService)

	server := server.NewServer(
		dbService,
		session,
		userRepo,
		listingRepo,
		userHandler,
		authHandler,
		listingHandler,
	)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	slog.Info("Server starting up...", slog.String("addr", server.Addr))

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	slog.Info("Graceful shutdown complete.")
}
