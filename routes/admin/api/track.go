package api

import (
	"encoding/json"
	"jinya-releases/database"
	"jinya-releases/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func getAllTracks(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]

	encoder := json.NewEncoder(w)
	tracks, err := database.Select[database.Track]("select * from track where application_id = $1", applicationId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not get tracks",
		})
		return
	}

	_ = encoder.Encode(tracks)
}

func getTrackById(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["id"]

	encoder := json.NewEncoder(w)
	track, err := database.Select[database.Track]("select * from track where application_id = $1 and id = $2", applicationId, trackId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not get track",
		})
		return
	}

	_ = encoder.Encode(track)
}

func createTrack(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	var track database.Track
	err := decoder.Decode(&track)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not decode track",
		})
		return
	}

	track.Id = uuid.New()
	track.ApplicationId = applicationId

	err = database.GetDbMap().Insert(&track)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not create track",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = encoder.Encode(track)
}

func updateTrack(w http.ResponseWriter, r *http.Request) {
	trackId := mux.Vars(r)["id"]

	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)

	var track database.Track
	err := decoder.Decode(&track)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not decode track",
		})
		return
	}

	dbTrack, err := database.Get[database.Track](trackId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not get track",
		})
		return
	}

	dbTrack.Name = track.Name
	dbTrack.Slug = track.Slug
	dbTrack.IsDefault = track.IsDefault

	_, err = database.GetDbMap().Update(dbTrack)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not update track",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteTrack(w http.ResponseWriter, r *http.Request) {
	applicationId := mux.Vars(r)["applicationId"]
	trackId := mux.Vars(r)["id"]

	encoder := json.NewEncoder(w)
	_, err := database.GetDbMap().Exec("delete from track where application_id = $1 and id = $2", applicationId, trackId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(utils.ErrorDetails{
			Message: "Could not delete track",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
