package handlers

import (
	"chi-test/internal/middleware"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"

	"chi-test/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(r chi.Router, db *gorm.DB) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
			signUpHandler(w, r, db)
		})
		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			loginHandler(w, r, db)
		})
	})
}

func signUpHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	log.Println("Signup handler called")
	var user models.User
	// Decode the request body into the User struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Println("Hashing user password...")
	// Hash the user's password before saving it to the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		http.Error(w, "Unable to create user", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	log.Println("Creating user in the database...")
	// Save the user to the database
	if err := db.Create(&user).Error; err != nil {
		log.Println("Error creating user in database:", err)
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	log.Println("User created successfully")
	w.WriteHeader(http.StatusCreated)
}

func loginHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	log.Println("Login handler called")
	var credentials models.User
	// Decode the request body into the User struct
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Println("Fetching user from database...")
	var user models.User
	// Fetch the user from the database by username
	if err := db.Where("username = ?", credentials.Username).First(&user).Error; err != nil {
		log.Println("User not found:", err)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	log.Println("Comparing passwords...")
	// Compare the hashed password with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		log.Println("Invalid credentials:", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	log.Println("Generating JWT token...")
	// Create a JWT token for the user
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime), // Set the token expiration time
		Subject:   fmt.Sprint(user.ID),                // Set the user ID as the subject of the token
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JwtKey)
	if err != nil {
		log.Println("Error signing token:", err)
		http.Error(w, "Unable to generate token", http.StatusInternalServerError)
		return
	}

	log.Println("Token generated successfully")
	// Return the token to the client
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
