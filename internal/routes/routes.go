package routes

import (
	"encoding/json"
	"net/http"

	"github.com/CVWO/sample-go-app/internal/handlers/users"
	"github.com/CVWO/sample-go-app/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func GetRoutes() func(r chi.Router) {
	return func(r chi.Router) {

		// public routes : routes that dont need to use the token aka the get stuff
		r.Post("/login", func(w http.ResponseWriter, req *http.Request) {
			response, _ := users.HandleUserLogin(w, req)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})
		r.Get("/posts/{postID}", func(w http.ResponseWriter, req *http.Request) {
				response, _ := users.HandleGetPost(w, req)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			})
		// login is post so data is more safe
		// routes that connect frontend to backend
		r.Post("/users", func(w http.ResponseWriter, req *http.Request) {
			response, _ := users.HandleCreateUser(w, req)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})

		r.Get("/posts", func(w http.ResponseWriter, req *http.Request) {
			response, _ := users.HandleGetPostsByCategory(w, req)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})
		r.Get("/comments", func(w http.ResponseWriter, req *http.Request) {
			response, _ := users.HandleGetComment(w, req)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})

		// private routes : routes that need the token
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)
			r.Post("/posts", func(w http.ResponseWriter, req *http.Request) {
				response, _ := users.HandleCreatePost(w, req)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			})

			r.Put("/posts", func(w http.ResponseWriter, req *http.Request) {
				response, _ := users.HandleLikesPost(w, req)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			})

			r.Post("/comments", func(w http.ResponseWriter, req *http.Request) {
				response, _ := users.HandleCreateComment(w, req)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			})

			r.Put("/comments", func(w http.ResponseWriter, req *http.Request) {
				response, _ := users.HandleLikesComment(w, req)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			})
			
			r.Get("/like/{postID}", func(w http.ResponseWriter, req *http.Request) {
				response, _ := users.HandleGetIndividualLike(w, req)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			})

		})

	}
}
