package main

import (
	customMiddleware "chi-test/internal/middleware"
	"log"
	"net/http"

	_ "chi-test/docs" // Import the generated Swagger docs
	"chi-test/internal/handlers"
	"chi-test/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Initialize the database
	db, err := database.Initialize()
	if err != nil {
		log.Fatal("failed to connect database")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)    // Middleware to log HTTP requests
	r.Use(middleware.Recoverer) // Middleware to recover from panics
	r.Use(customMiddleware.JwtAuthentication)

	handlers.RegisterAuthRoutes(r, db)
	handlers.RegisterStaticRoutes(r)
	handlers.RegisterProtectedRoutes(r, db)

	// Swagger documentation endpoint
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Println("Starting server on :8080")
	// Start the HTTP server on port 8080
	http.ListenAndServe(":8080", r)
}
