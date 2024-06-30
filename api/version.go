package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"jinya-releases/service"
	"log"
	"net/http"
)

func getAllVersions(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["trackId"]
	encoder := json.NewEncoder(w)
	versions, errDetails, errStatus := service.GetAllVersions(applicationId, trackId)
	if errDetails != nil {
		w.WriteHeader(errStatus)
		_ = encoder.Encode(errDetails)
		return
	}

	_ = encoder.Encode(versions)
}

func getVersionById(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["trackId"]
	versionId := mux.Vars(r)["id"]
	encoder := json.NewEncoder(w)
	version, errDetails, errStatus := service.GetVersionById(applicationId, trackId, versionId)
	if errDetails != nil {
		w.WriteHeader(errStatus)
		_ = encoder.Encode(errDetails)
		return
	}

	_ = encoder.Encode(version)
}

func createVersion(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["trackId"]
	encoder := json.NewEncoder(w)
	version, errDetails, errStatus := service.CreateVersion(r.Body, applicationId, trackId)
	if errDetails != nil {
		w.WriteHeader(errStatus)
		_ = encoder.Encode(errDetails)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = encoder.Encode(version)
}

func deleteVersion(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["trackId"]
	versionId := mux.Vars(r)["id"]
	encoder := json.NewEncoder(w)
	errDetails, status := service.DeleteVersion(applicationId, trackId, versionId)
	if errDetails != nil {
		w.WriteHeader(status)
		_ = encoder.Encode(errDetails)
		return
	}

	w.WriteHeader(status)
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

func pushVersion(w http.ResponseWriter, r *http.Request) {
	applicationSlug := mux.Vars(r)["applicationSlug"]
	trackSlug := mux.Vars(r)["trackSlug"]
	versionNumber := mux.Vars(r)["version"]
	encoder := json.NewEncoder(w)
	errDetails, status := service.PushVersion(r, applicationSlug, trackSlug, versionNumber)
	if errDetails != nil {
		w.WriteHeader(status)
		_ = encoder.Encode(errDetails)
		return
	}

	w.WriteHeader(status)
}

func uploadVersion(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["trackId"]
	versionNumber := mux.Vars(r)["version"]
	encoder := json.NewEncoder(w)
	errDetails, status := service.UploadVersion(r, applicationId, trackId, versionNumber)
	if errDetails != nil {
		w.WriteHeader(status)
		_ = encoder.Encode(errDetails)
		return
	}

	w.WriteHeader(status)
}
