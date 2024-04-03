package api

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type createTrackRequest struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	IsDefault bool   `json:"isDefault"`
}

type updateTrackRequest struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	IsDefault bool   `json:"isDefault"`
}

func getAllTracks(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	log.Printf("App id: %s", applicationId)

	w.WriteHeader(http.StatusNotImplemented)
}

func getTrackById(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["id"]
	log.Printf("App id: %s", applicationId)
	log.Printf("Track id: %s", trackId)

	w.WriteHeader(http.StatusNotImplemented)
}

func createTrack(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	log.Printf("App id: %s", applicationId)

	w.WriteHeader(http.StatusNotImplemented)
}

func updateTrack(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["id"]
	log.Printf("App id: %s", applicationId)
	log.Printf("Track id: %s", trackId)

	w.WriteHeader(http.StatusNotImplemented)
}

func deleteTrack(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["id"]
	log.Printf("App id: %s", applicationId)
	log.Printf("Track id: %s", trackId)

	w.WriteHeader(http.StatusNotImplemented)
}
