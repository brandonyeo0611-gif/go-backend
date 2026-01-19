package router

import (
	"github.com/CVWO/sample-go-app/internal/database"
	"github.com/CVWO/sample-go-app/internal/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)


func Setup(db *database.Database) chi.Router {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000","https://react-1-qgqh.onrender.com",}, // your frontend URL
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        AllowCredentials: true,
    }))
	// connecting backend to frontend!! stuck here for ages
	// if not cannot connect
	setUpRoutes(r, db)
	return r
}

func setUpRoutes(r chi.Router, db *database.Database) {
	r.Group(routes.GetRoutes(db))
}
