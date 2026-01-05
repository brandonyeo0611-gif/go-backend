package routes

import (
	"encoding/json"
	"net/http"

	"github.com/CVWO/sample-go-app/internal/database"
	"github.com/CVWO/sample-go-app/internal/handlers/users"
	"github.com/CVWO/sample-go-app/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func GetRoutes(db *database.Database) func(r chi.Router) {
	return func(r chi.Router) {

		// public routes : routes that dont need to use the token aka the get stuff
		r.Post("/login", func(w http.ResponseWriter, req *http.Request) {
			response, _ := users.HandleUserLogin(w, req, db)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})
		r.Get("/posts/{postID}", func(w http.ResponseWriter, req *http.Request) {
			response, _ := users.HandleGetPost(w, req, db)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})
		// login is post so data is more safe
		// routes that connect frontend to backend
		r.Post("/users", func(w http.ResponseWriter, req *http.Request) {
			response, _ := users.HandleCreateUser(w, req, db)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})

		r.Get("/posts", func(w http.ResponseWriter, req *http.Request) {
			response, _ := users.HandleGetPostsByCategory(w, req, db)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})
		r.Get("/comments", func(w http.ResponseWriter, req *http.Request) {
			response, _ := users.HandleGetComment(w, req, db)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})
		r.Post("/RefreshAccessToken", func(w http.ResponseWriter, req *http.Request) {
			response, _ := users.HandleRefreshAccessToken(w, req, db)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})

		// private routes : routes that need the token
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)
			r.Post("/posts", func(w http.ResponseWriter, req *http.Request) {
				response, _ := users.HandleCreatePost(w, req, db)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			})

			r.Put("/posts", func(w http.ResponseWriter, req *http.Request) {
				response, _ := users.HandleLikesPost(w, req, db)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			})

			r.Post("/comments", func(w http.ResponseWriter, req *http.Request) {
				response, _ := users.HandleCreateComment(w, req, db)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			})

			r.Put("/comments", func(w http.ResponseWriter, req *http.Request) {
				response, _ := users.HandleLikesComment(w, req, db)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			})

			r.Get("/like/{postID}", func(w http.ResponseWriter, req *http.Request) {
				response, _ := users.HandleGetIndividualLike(w, req, db)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			})

			r.Put("/users/profile_pic", func(w http.ResponseWriter, req *http.Request) {
				response, _ := users.HandleChangeProfilePic(w, req, db)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			})

			r.Get("/users/profile_pic", func(w http.ResponseWriter, req *http.Request) {
				response, _ := users.HandleGetProfilePic(w, req, db)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			})

		})

	}
}
