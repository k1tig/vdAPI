package main

import (
	"encoding/json"
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
	/*router.HandleFunc("GET /racer/{raceID}", func(w http.ResponseWriter, r *http.Request) {
		raceID := r.PathValue("raceID")
		w.Write([]byte("Race ID: " + raceID))
	})*/
	router.HandleFunc("GET /racer", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("\n\nRacers: "))
		for _, i := range racers {
			w.Write([]byte(i.Name + " "))
		}
	})
	router.HandleFunc("POST /racer", func(w http.ResponseWriter, r *http.Request) {
		var racer racer
		err := json.NewDecoder(r.Body).Decode(&racer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		racers = append(racers, racer)
		w.Write([]byte("\n\nEntered: " + racer.Name))

	})

	server := http.Server{
		Addr:    s.addr,
		Handler: middleware.Logging(router),
	}
	log.Println("server starting on:", s.addr)
	return server.ListenAndServe()
}

// curl -X POST -H "Content-Type: application/json" -d '{"racername":"MeeDok"}' http://localhost:8080/racer
