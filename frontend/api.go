package frontend

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/iancoleman/orderedmap"
	"jinya-releases/database/models"
	"net/http"
)

func versionsToMap(versions []models.Version) *orderedmap.OrderedMap {
	versionMap := orderedmap.New()
	for _, version := range versions {
		versionMap.Set(version.Version, version.Url)
	}

	return versionMap
}

func getApplicationJson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	applicationSlug := vars["applicationSlug"]

	app, err := models.GetApplicationBySlug(applicationSlug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tracks, err := models.GetAllTracks(app.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)

	track := (*models.Track)(nil)

	for _, t := range tracks {
		if t.IsDefault {
			track = &t
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if track == nil {
		if len(tracks) >= 1 {
			track = &tracks[0]
		} else if len(tracks) == 0 {
			_ = encoder.Encode([]any{})
			return
		}
	}

	versions, err := models.GetAllVersions(app.Id, tracks[0].Id)
	if err != nil {
		versions = []models.Version{}
	}

	_ = encoder.Encode(versionsToMap(versions))
}

func getTrackJson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	applicationSlug := vars["applicationSlug"]
	trackSlug := vars["trackSlug"]

	encoder := json.NewEncoder(w)

	w.Header().Set("Content-Type", "application/json")
	versions, err := models.GetVersionBySlugs(applicationSlug, trackSlug)
	if err != nil {
		versions = []models.Version{}
	}

	_ = encoder.Encode(versionsToMap(versions))
}
