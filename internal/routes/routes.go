package routes

import (
	"encoding/json"
	"net/http"

	"github.com/CVWO/sample-go-app/internal/handlers/users"
	"github.com/go-chi/chi/v5"
)

func GetRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/login", func(w http.ResponseWriter, req *http.Request) {
			response, _ := users.HandleUserLogin(w, req)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})
		// login is post so data is more safe

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

		r.Get("/comments", func(w http.ResponseWriter, req *http.Request) {
			response, _ := users.HandleGetComment(w, req)
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

	}
}
