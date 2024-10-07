package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RegisterStaticRoutes(r chi.Router) {
	r.Route("/static", func(r chi.Router) {
		r.Handle("/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/assets/"))))
	})
}
