package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"jinya-releases/service"
	"net/http"
)

func createPushtoken(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	pushtoken, errDetails, errStatus := service.CreatePushtoken(r.Body)
	if errDetails != nil {
		w.WriteHeader(errStatus)
		_ = encoder.Encode(errDetails)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = encoder.Encode(pushtoken)
}

func getAllPushtokens(w http.ResponseWriter, _ *http.Request) {
	encoder := json.NewEncoder(w)
	pushtokens, errDetails := service.GetAllPushtokens()
	if errDetails != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(errDetails)
		return
	}

	_ = encoder.Encode(pushtokens)
}

func getPushtokenById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	encoder := json.NewEncoder(w)
	pushtoken, errDetails, errStatus := service.GetPushtokenById(id)
	if errDetails != nil {
		w.WriteHeader(errStatus)
		_ = encoder.Encode(errDetails)
		return
	}

	_ = encoder.Encode(pushtoken)
}

func updatePushtoken(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	encoder := json.NewEncoder(w)
	_, errDetails, errStatus := service.UpdatePushtoken(id, r.Body)
	if errDetails != nil {
		w.WriteHeader(errStatus)
		_ = encoder.Encode(errDetails)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deletePushtoken(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	encoder := json.NewEncoder(w)
	errDetails, status := service.DeletePushtoken(id)
	if errDetails != nil {
		w.WriteHeader(status)
		_ = encoder.Encode(errDetails)
		return
	}

	w.WriteHeader(status)
}
