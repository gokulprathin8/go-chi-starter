package handlers

import (
	"chi-test/internal/middleware"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func RegisterProtectedRoutes(r chi.Router, db *gorm.DB) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuthentication)
		r.Get("/protected", protectedHandler)
	})
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Protected handler called")
	// Respond with a message indicating the endpoint is protected
	w.Write([]byte("This is a protected endpoint"))
}
