package api

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type createApplicationRequest struct {
	Name                 string `json:"name"`
	Slug                 string `json:"slug"`
	HomepageTemplate     string `json:"homepageTemplate"`
	TrackpageTemplate    string `json:"trackpageTemplate"`
	AdditionalCss        string `json:"additionalCss,omitempty"`
	AdditionalJavaScript string `json:"additionalJavaScript,omitempty"`
}

type updateApplicationRequest struct {
	Name                 string `json:"name"`
	Slug                 string `json:"slug"`
	HomepageTemplate     string `json:"homepageTemplate"`
	TrackpageTemplate    string `json:"trackpageTemplate"`
	AdditionalCss        string `json:"additionalCss,omitempty"`
	AdditionalJavaScript string `json:"additionalJavaScript,omitempty"`
}

func getAllApplications(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func getApplicationById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("Id: %s", id)

	w.WriteHeader(http.StatusNotImplemented)
}

func createApplication(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func updateApplication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("Id: %s", id)

	w.WriteHeader(http.StatusNotImplemented)
}

func deleteApplication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("Id: %s", id)

	w.WriteHeader(http.StatusNotImplemented)
}

func uploadLogo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("Id: %s", id)

	w.WriteHeader(http.StatusNotImplemented)
}
