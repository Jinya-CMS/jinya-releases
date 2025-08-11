package api

import (
	"encoding/json"
	"jinya-releases/service"
	"net/http"

	"github.com/gorilla/mux"
)

func getAllTracks(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]

	encoder := json.NewEncoder(w)
	tracks, errDetails, errStatus := service.GetAllTracks(applicationId)
	if errDetails != nil {
		w.WriteHeader(errStatus)
		_ = encoder.Encode(errDetails)
		return
	}

	_ = encoder.Encode(tracks)
}

func getTrackById(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["id"]

	encoder := json.NewEncoder(w)
	track, errDetails, errStatus := service.GetTrackById(trackId, applicationId)
	if errDetails != nil {
		w.WriteHeader(errStatus)
		_ = encoder.Encode(errDetails)
		return
	}

	_ = encoder.Encode(track)
}

func createTrack(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	encoder := json.NewEncoder(w)
	track, errDetails, errStatus := service.CreateTrack(r.Body, applicationId)
	if errDetails != nil {
		w.WriteHeader(errStatus)
		_ = encoder.Encode(errDetails)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = encoder.Encode(track)
}

func updateTrack(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["id"]

	encoder := json.NewEncoder(w)
	_, errDetails, errStatus := service.UpdateTrack(trackId, applicationId, r.Body)
	if errDetails != nil {
		w.WriteHeader(errStatus)
		_ = encoder.Encode(errDetails)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteTrack(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["id"]

	encoder := json.NewEncoder(w)
	errDetails, status := service.DeleteTrack(trackId, applicationId)
	if errDetails != nil {
		w.WriteHeader(status)
		_ = encoder.Encode(errDetails)
		return
	}

	w.WriteHeader(status)
}
