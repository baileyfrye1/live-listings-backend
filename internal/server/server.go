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

	db             database.Service
	session        *session.Session
	userRepo       *repo.UserRepository
	listingRepo    *repo.ListingRepository
	userHandler    *handler.UserHandler
	authHandler    *handler.AuthHandler
	listingHandler *handler.ListingHandler
}

func NewServer(
	db database.Service,
	session *session.Session,
	userRepo *repo.UserRepository,
	listingRepo *repo.ListingRepository,
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	listingHandler *handler.ListingHandler,
) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:           port,
		db:             db,
		session:        session,
		userRepo:       userRepo,
		listingRepo:    listingRepo,
		userHandler:    userHandler,
		authHandler:    authHandler,
		listingHandler: listingHandler,
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
