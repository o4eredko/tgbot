package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/sheets/v4"

	sheets2 "tg-bot/internal/sheets"
)

func main() {
	privateKey, err := os.ReadFile("./private.key")
	if err != nil {
		log.Panic(fmt.Errorf("read private key from file: %w", err))
	}

	conf := &jwt.Config{
		Email:        "telegram@folderly-app.iam.gserviceaccount.com",
		PrivateKey:   privateKey,
		PrivateKeyID: "e571adfd1caa0727af76f7abd9c5a827f5ed9371",
		TokenURL:     "https://oauth2.googleapis.com/token",
		Scopes: []string{
			"https://www.googleapis.com/auth/spreadsheets.readonly",
		},
	}
	spreadsheetID := "1jUKFYm03vwC67P9mbKhGV8icIMgy5vLNMaCS9EQ0kWk"

	client := conf.Client(context.Background())

	// Create a service object for Google sheets
	sheetsClient, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	service := sheets2.NewService(spreadsheetID, sheetsClient)
	router := mux.NewRouter()
	handler := NewHandler(service)
	handler.RegisterRoutes(router)

	listen := "0.0.0.0:8080"
	log.Println("Starting listening on:", listen)

	if err := http.ListenAndServe(listen, router); err != nil {
		log.Panic(fmt.Sprintf("listen http: %s", err))
	}
}
