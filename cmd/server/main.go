package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type messageRequest struct {
	Message string `json:"message"`
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Heartbeat("/health"))

	router.Post("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		var message messageRequest
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("A mensagem que chegou via request foi: %s\n", message.Message)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		if err := json.NewEncoder(w).Encode("data"); err != nil {
			log.Fatalf(err.Error())
		}
	})

	server := http.Server{
		ReadTimeout:       time.Duration(10) * time.Second,
		ReadHeaderTimeout: time.Duration(10) * time.Second,
		Handler:           router,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", "3000"))
	if err != nil {
		panic(err)
	}
	server.Serve(listener)
}
