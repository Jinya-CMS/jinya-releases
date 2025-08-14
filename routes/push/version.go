package push

import (
	"encoding/json"
	"jinya-releases/database"
	"jinya-releases/storage"
	"jinya-releases/utils"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func pushVersion(w http.ResponseWriter, r *http.Request) {
	applicationSlug := mux.Vars(r)["applicationSlug"]
	trackSlug := mux.Vars(r)["trackSlug"]
	versionNumber := mux.Vars(r)["versionNumber"]
	encoder := json.NewEncoder(w)

	track, err := database.SelectOne[database.Track]("select t.* from track t inner join application a on a.id = t.application_id and a.slug = $1 where t.slug = $2", applicationSlug, trackSlug)
	if err != nil {
		log.Printf("Could not find track: %v", err)
		w.WriteHeader(http.StatusNotFound)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not find track",
		})
		return
	}

	err = storage.UploadVersion(r.Body, r.ContentLength, r.Header.Get("Content-Type"), &database.Version{
		Id:            uuid.UUID{},
		ApplicationId: track.ApplicationId,
		TrackId:       track.Id.String(),
		Version:       versionNumber,
		UploadDate:    time.Time{},
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
