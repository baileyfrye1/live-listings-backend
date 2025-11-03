package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	cm "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"server/internal/server/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	authMiddleware := middleware.Authenticate(s.session, s.userRepo)
	authorizeMiddleware := middleware.Authorize()

	r.Use(cm.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/listings", s.listingHandler.GetAllListings)
		r.Get("/listings/{listingId}", s.listingHandler.GetListingById)
		r.Patch("/listings/{listingId}/views", s.listingHandler.TrackViewsByListingId)

		r.Get("/agents", s.userHandler.GetAllAgents)
		r.Get("/agents/{agentId}", s.userHandler.GetAgentById)
		r.Get("/agents/{agentId}/listings", s.listingHandler.GetAgentListings)

		r.Route("/auth", func(u chi.Router) {
			u.Post("/register", s.authHandler.Register)
			u.Post("/login", s.authHandler.Login)
		})
	})

	// Authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Get("/users/profile", s.userHandler.GetCurrentUser)
		r.Patch("/users/profile", s.userHandler.UpdateUserById)
		r.Post("/auth/logout", s.authHandler.Logout)

		r.Get("/favorites", s.favoriteHandler.GetUserFavorites)
		r.Post("/favorites", s.favoriteHandler.CreateFavorite)
		r.Delete("/favorites/{listingId}", s.favoriteHandler.DeleteFavoriteByListingId)

		r.Get("/notifications", s.notificationHandler.GetAllNotificationsByUserId)
		r.Post("/notifications", s.notificationHandler.CreateNotification)
		r.Patch(
			"/notifications/{notificationId}",
			s.notificationHandler.ToggleNotificationReadStatus,
		)

		r.Get("/ws", s.wsManager.StartWSConn)

		// Agent/admin routes
		r.Group(func(r chi.Router) {
			r.Use(authorizeMiddleware)

			r.Get("/agents/me/listings", s.listingHandler.GetMyListings)
			r.Post("/listings", s.listingHandler.CreateListing)
			r.Patch("/listings/{listingId}", s.listingHandler.UpdateMyListing)
			r.Delete("/listings/{listingId}", s.listingHandler.DeleteMyListing)

			r.Get("/users", s.userHandler.GetAllUsers)
			r.Patch("/users/{userId}", s.userHandler.UpdateUserById)
		})
	})

	r.Get("/health", s.healthHandler)

	return r
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
