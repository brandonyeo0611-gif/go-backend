package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CVWO/sample-go-app/internal/router"
	"github.com/CVWO/sample-go-app/internal/database"
)

func main() {
	db, err := database.GetDB()
	if err != nil {
		log.Fatalln("Failed to connect to DB:", err)
	}
	r := router.Setup(db)
	fmt.Print("Listening on port 8000 at http://localhost:8000!")

	log.Fatalln(http.ListenAndServe(":8000", r))
}
