package api

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type createVersionRequest struct {
	Version string `json:"version"`
}

func getAllVersions(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["trackId"]
	log.Printf("App id: %s", applicationId)
	log.Printf("Track id: %s", trackId)

	w.WriteHeader(http.StatusNotImplemented)
}

func getVersionById(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["trackId"]
	versionId := mux.Vars(r)["id"]
	log.Printf("App id: %s", applicationId)
	log.Printf("Track id: %s", trackId)
	log.Printf("Version id: %s", versionId)

	w.WriteHeader(http.StatusNotImplemented)
}

func createVersion(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["trackId"]
	log.Printf("App id: %s", applicationId)
	log.Printf("Track id: %s", trackId)

	w.WriteHeader(http.StatusNotImplemented)
}

func deleteVersion(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["trackId"]
	versionId := mux.Vars(r)["id"]
	log.Printf("App id: %s", applicationId)
	log.Printf("Track id: %s", trackId)
	log.Printf("Version id: %s", versionId)

	w.WriteHeader(http.StatusNotImplemented)
}

func uploadVersionBinary(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["trackId"]
	versionId := mux.Vars(r)["id"]
	log.Printf("App id: %s", applicationId)
	log.Printf("Track id: %s", trackId)
	log.Printf("Version id: %s", versionId)

	w.WriteHeader(http.StatusNotImplemented)
}
