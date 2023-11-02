package main

import (
	"fmt"
	"log"
	"net/http"

	"goAlertStore/data"

	"goAlertStore/api"

	"github.com/go-chi/chi"
)

func main() {
	// Initializing the data layer
	db := data.NewDatabase("alert1.db")
	if err := db.Open(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initializing API handlers
	alertHandler := api.NewAlertHandler(db)

	// Router setup
	r := chi.NewRouter()
	api.RegisterRoutes(r, alertHandler)

	// Server start
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", "8080"),
		Handler: r,
	}
	log.Printf("Server started... %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf(fmt.Sprintf("%+v", err))
	}
}
