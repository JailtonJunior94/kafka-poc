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

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Heartbeat("/health"))

	router.Post("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		log.Println("chjamou")

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

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", "5000"))
	if err != nil {
		panic(err)
	}
	server.Serve(listener)
}
