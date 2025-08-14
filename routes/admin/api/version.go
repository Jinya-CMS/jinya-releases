package api

import (
	"encoding/json"
	"jinya-releases/config"
	"jinya-releases/database"
	"jinya-releases/storage"
	"jinya-releases/utils"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const baseVersionQuery = `
select v.*, $1 || '/content/version/' || a.slug || '/' || t.slug || '/' || v.version as url
from version v
         inner join application a on a.id = v.application_id
         inner join track t on t.id = v.track_id
where v.application_id = $2
  and v.track_id = $3
`

func getAllVersions(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["trackId"]
	encoder := json.NewEncoder(w)
	versions, err := database.Select[database.Version](baseVersionQuery, config.LoadedConfiguration.ServerUrl, applicationId, trackId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not get versions",
		})
		return
	}

	_ = encoder.Encode(versions)
}

func getVersionById(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["trackId"]
	versionId := mux.Vars(r)["id"]
	encoder := json.NewEncoder(w)
	version, err := database.Select[database.Version](baseVersionQuery+" and id = $4", config.LoadedConfiguration.ServerUrl, applicationId, trackId, versionId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not get version",
		})
		return
	}

	_ = encoder.Encode(version)
}

func deleteVersion(w http.ResponseWriter, r *http.Request) {
	versionId := mux.Vars(r)["id"]
	encoder := json.NewEncoder(w)

	_, err := database.GetDbMap().Exec("delete from version where id = $1", versionId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not delete version",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func uploadVersion(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["trackId"]
	versionNumber := mux.Vars(r)["version"]
	encoder := json.NewEncoder(w)
	err := storage.UploadVersion(r.Body, r.ContentLength, r.Header.Get("Content-Type"), &database.Version{
		Id:            uuid.UUID{},
		ApplicationId: applicationId,
		TrackId:       trackId,
		Version:       versionNumber,
	})

	if err != nil {
		log.Printf("Could not upload version: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not upload version",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
}
