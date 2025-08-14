package api

import (
	"encoding/json"
	"jinya-releases/database"
	"jinya-releases/storage"
	"jinya-releases/utils"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func getAllApplications(w http.ResponseWriter, _ *http.Request) {
	encoder := json.NewEncoder(w)
	applications, err := database.Select[database.Application]("select * from application")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not get applications",
		})
		return
	}

	_ = encoder.Encode(applications)
}

func getApplicationById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	encoder := json.NewEncoder(w)

	application, err := database.Get[database.Application](id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not get application",
		})
		return
	}

	_ = encoder.Encode(application)
}

func createApplication(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	var app database.Application
	err := decoder.Decode(&app)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not decode application",
		})
		return
	}

	app.Id = uuid.New()

	err = database.GetDbMap().Insert(&app)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not create application",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = encoder.Encode(app)
}

func updateApplication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	var app database.Application
	err := decoder.Decode(&app)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not decode application",
		})
		return
	}

	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not parse application id",
		})
		return
	}

	app.Id = parsedId

	_, err = database.GetDbMap().Update(app)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not update application",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteApplication(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	encoder := json.NewEncoder(w)
	_, err := database.GetDbMap().Exec("delete from application where id = $1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not delete application",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func resetToken(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	encoder := json.NewEncoder(w)
	_, err := database.GetDbMap().Exec("delete from push_token where application_id = $1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not delete tokens",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func createToken(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	encoder := json.NewEncoder(w)
	token := database.PushToken{
		Id:            uuid.New(),
		Token:         uuid.NewString(),
		ApplicationId: id,
	}
	err := database.GetDbMap().Insert(&token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not create token",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = encoder.Encode(token)
}

func uploadLogo(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	err := storage.UploadLogo(r)
	if err != nil {
		log.Printf("Failed to upload logo: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not upload logo",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
