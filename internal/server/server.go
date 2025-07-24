package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"server/database"
	"server/internal/api/handler"
	"server/internal/repo"
	"server/internal/session"
)

type Server struct {
	port int

	db          database.Service
	session     *session.Session
	userRepo    *repo.UserRepository
	userHandler *handler.UserHandler
	authHandler *handler.AuthHandler
}

func NewServer(
	db database.Service,
	session *session.Session,
	userRepo *repo.UserRepository,
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:        port,
		db:          db,
		session:     session,
		userRepo:    userRepo,
		userHandler: userHandler,
		authHandler: authHandler,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
