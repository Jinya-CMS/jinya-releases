package api

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type createTokenRequest struct {
	AllowedApps   []string `json:"allowedApps,omitempty"`
	AllowedTracks []string `json:"allowedTracks,omitempty"`
}

type updateTokenRequest struct {
	AllowedApps   []string `json:"allowedApps,omitempty"`
	AllowedTracks []string `json:"allowedTracks,omitempty"`
}

func getAllPushTokens(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("Id: %s", id)

	w.WriteHeader(http.StatusNotImplemented)
}

func getPushTokenById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("Id: %s", id)

	w.WriteHeader(http.StatusNotImplemented)
}

func createPushToken(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func updatePushToken(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("Id: %s", id)

	w.WriteHeader(http.StatusNotImplemented)
}

func deletePushToken(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("Id: %s", id)

	w.WriteHeader(http.StatusNotImplemented)
}
