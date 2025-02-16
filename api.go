package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	router.HandleFunc("GET /groups", getGroups)
	router.HandleFunc("GET /group/:id", getGroupById)
	router.HandleFunc("POST /groups", createGroup)

	server := http.Server{
		Addr:    s.addr,
		Handler: middleware.Logging(router),
	}
	log.Println("server starting on:", s.addr)
	return server.ListenAndServe()
}

func getGroups(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(racerGroups)
	if err != nil {
		http.Error(w, "Error Marshalling JSON", http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

func getGroupById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
	}
	for _, i := range racerGroups {
		if i.GroupId == id {
			resp, err := json.Marshal(i)
			if err != nil {
				http.Error(w, "Error Marshalling JSON", http.StatusInternalServerError)
				return
			}
			w.Write(resp)
		}
	}

}

func createGroup(w http.ResponseWriter, r *http.Request) {

	//probably need media type verification?
	dec := json.NewDecoder(r.Body)

	var rg raceGroup
	err := dec.Decode(&rg)
	// just a bookmark for future incoming data handling
	//dec.DisallowUnknownFields()

	if err != nil {
		http.Error(w, "Error Marshalling JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Group Created Successfully"}`))
}

func updateGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
	}

	dec := json.NewDecoder(r.Body)
	var rgUpdate raceGroup
	err = dec.Decode(&rgUpdate)

	if err != nil {
		http.Error(w, "Error Marshalling JSON", http.StatusInternalServerError)
		return
	}

	for _, i := range racerGroups {
		if i.GroupId == id {
			//do stuff
		}
	}

}

// The group.rev PUT needs to be managed in a way to not acccept certain versions.

// curl -X POST -H "Content-Type: application/json" -d '{"racername":"MeeDok"}' http://localhost:8080/racer

// curl -X GET http://localhost:8080/racer
