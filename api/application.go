package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"jinya-releases/service"
	"jinya-releases/storage"
	"net/http"
)

func getAllApplications(w http.ResponseWriter, _ *http.Request) {
	encoder := json.NewEncoder(w)
	applications, errDetails := service.GetAllApplications()
	if errDetails != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(errDetails)
		return
	}

	_ = encoder.Encode(applications)
}

func getApplicationById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	encoder := json.NewEncoder(w)
	application, errDetails, errStatus := service.GetApplicationById(id)
	if errDetails != nil {
		w.WriteHeader(errStatus)
		_ = encoder.Encode(errDetails)
		return
	}

	_ = encoder.Encode(application)
}

func createApplication(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	app, errDetails, errStatus := service.CreateApplication(r.Body)
	if errDetails != nil {
		w.WriteHeader(errStatus)
		_ = encoder.Encode(errDetails)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = encoder.Encode(app)
}

func updateApplication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	encoder := json.NewEncoder(w)
	_, errDetails, errStatus := service.UpdateApplication(id, r.Body)
	if errDetails != nil {
		w.WriteHeader(errStatus)
		_ = encoder.Encode(errDetails)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteApplication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	encoder := json.NewEncoder(w)
	errDetails, status := service.DeleteApplication(id)
	if errDetails != nil {
		w.WriteHeader(status)
		_ = encoder.Encode(errDetails)
		return
	}

	w.WriteHeader(status)
}

func resetToken(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	encoder := json.NewEncoder(w)
	errDetails, status := service.ResetToken(id)
	if errDetails != nil {
		w.WriteHeader(status)
		_ = encoder.Encode(errDetails)
		return
	}

	w.WriteHeader(status)
}

func createToken(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	encoder := json.NewEncoder(w)
	errDetails, token, status := service.CreatePushToken(id)
	if errDetails != nil {
		w.WriteHeader(status)
		_ = encoder.Encode(errDetails)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = encoder.Encode(token)
}

func uploadLogo(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	errDetails, status := storage.UploadLogo(r)
	if errDetails != nil {
		w.WriteHeader(status)
		_ = encoder.Encode(errDetails)
		return
	}

	w.WriteHeader(status)
}
