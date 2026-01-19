package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CVWO/sample-go-app/internal/database"
	"github.com/CVWO/sample-go-app/internal/router"
)

func main() {
	db, err := database.GetDB()
	if err != nil {
		log.Fatalln("Failed to connect to DB:", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	r := router.Setup(db)
	fmt.Print("Listening on port 8000 at https://brandonwebforumgobackend.onrender.com!")

	log.Fatalln(http.ListenAndServe(":"+port, r))
}
