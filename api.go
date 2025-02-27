package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/k1tig/vdAPI/middleware"
)

type APIserver struct {
	addr string
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"` // interface spot?
	Message string      `json:"message"`
}

var mutex sync.Mutex

func NewAPIServer(addr string) *APIserver {
	return &APIserver{
		addr: addr,
	}
}

func (s *APIserver) Run() error {
	router := http.NewServeMux() //list routes below

	//runs top to bottom
	stack := middleware.CreateStack(
		middleware.StripTrailingSlash,
		middleware.Logging,
		middleware.ApiKeyMiddleware,
	)

	router.HandleFunc("GET /groups", getGroups)
	router.HandleFunc("GET /groups/{id}", getGroupById)
	router.HandleFunc("POST /groups", createGroup)
	router.HandleFunc("PUT /groups/{id}", updateGroup)
	router.HandleFunc("DELETE /groups/{id}", deleteGroup)

	server := http.Server{
		Addr:    s.addr,
		Handler: stack(router),
	}
	log.Println("server starting on:", s.addr)
	return server.ListenAndServe()
}

func getGroups(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	data := racerGroupResponse{
		Groups: racerGroups,
	}
	respStruct := APIResponse{
		Success: true,
		Data:    data,
	}
	resp := makeResponse(respStruct, w)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func getGroupById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
	}
	mutex.Lock()
	defer mutex.Unlock()
	for _, group := range racerGroups {
		if group.GroupId == id {
			resp, err := json.Marshal(group)
			if err != nil {
				http.Error(w, "Error Marshalling JSON", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`"Message" : "Group not found"`))
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
	if rg.GroupPhrase == "" {
		respStruct := APIResponse{
			Success: false,
			Message: "Passphrase void",
		}
		resp := makeResponse(respStruct, w)
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		return
	}
	// Not sure if a more specific lock is justified?
	rg.GroupId = groupCounter
	groupCounter++
	currentTime := time.Now()
	rg.Lifetime = currentTime

	mutex.Lock()
	racerGroups = append(racerGroups, rg)
	mutex.Unlock()

	respStruct := APIResponse{
		Success: true,
		Message: "Group Created Succeffully",
	}

	resp := makeResponse(respStruct, w)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
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

	mutex.Lock()
	defer mutex.Unlock()
	for targetGroup, i := range racerGroups {
		if i.GroupId == id {
			if rgUpdate.GroupPhrase != i.GroupPhrase {
				respStruct := APIResponse{
					Success: false,
					Message: "Permision denied: Passphrase incorrect",
				}
				resp := makeResponse(respStruct, w)
				w.WriteHeader(http.StatusOK)
				w.Write(resp)
				break
			}
			if i.GroupRev != rgUpdate.GroupRev {
				respStruct := APIResponse{
					Success: false,
					Message: "Error: Client group revision requires update",
				}
				resp := makeResponse(respStruct, w)
				w.WriteHeader(http.StatusOK)
				w.Write(resp)
				break
			}

			//added work to replace whole bracket instead of just updating, figure out better way later
			racerGroups[targetGroup] = rgUpdate
			racerGroups[targetGroup].GroupRev++
			racerGroups[targetGroup].Lifetime = time.Now()
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Group Updated Successfully"}`))
	}
}

func deleteGroup(w http.ResponseWriter, r *http.Request) {
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

	mutex.Lock()
	defer mutex.Unlock()
	for targetGroup, i := range racerGroups {
		if i.GroupId == id {
			if rgUpdate.GroupPhrase != i.GroupPhrase {
				respStruct := APIResponse{
					Success: false,
					Message: "Permision denied: Passphrase incorrect",
				}
				resp := makeResponse(respStruct, w)
				w.WriteHeader(http.StatusOK)
				w.Write(resp)
				break
			}

			//added work to replace whole bracket instead of just updating, figure out better way later
			racerGroups = append(racerGroups[:targetGroup], racerGroups[targetGroup+1:]...)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Group Updated Successfully"}`))
	}

}

// Need to standardize the JSON response format
func makeResponse(response APIResponse, w http.ResponseWriter) []byte {
	respStruct := response
	resp, err := json.Marshal(respStruct)
	if err != nil {
		http.Error(w, "Error Marshalling JSON", http.StatusInternalServerError)
		return nil
	}
	return resp
}
