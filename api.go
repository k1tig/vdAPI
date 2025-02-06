package main

import (
	"log"
	"net/http"

	"github.com/k1tig/vdAPI/middleware"
)

type APIserver struct {
	addr string
}

func NewAPIServer(addr string) *APIserver {
	return &APIserver{
		addr: addr,
	}
}

func (s *APIserver) Run() error {
	router := http.NewServeMux() //list routes below
	router.HandleFunc("GET /races/{raceID}", func(w http.ResponseWriter, r *http.Request) {
		raceID := r.PathValue("raceID")
		w.Write([]byte("Race ID: " + raceID))
	})

	server := http.Server{
		Addr:    s.addr,
		Handler: middleware.Logging(router),
	}
	log.Println("server starting on:", s.addr)
	return server.ListenAndServe()
}
