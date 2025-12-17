package router

import (
	"github.com/CVWO/sample-go-app/internal/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)


func Setup() chi.Router {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"}, // your frontend URL
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        AllowCredentials: true,
    }))
	setUpRoutes(r)
	return r
}

func setUpRoutes(r chi.Router) {
	r.Group(routes.GetRoutes())
}
