package content

import (
	"github.com/gorilla/mux"
	"io"
	"jinya-releases/database/models"
	"jinya-releases/storage"
	"net/http"
)

func GetVersion(w http.ResponseWriter, r *http.Request) {
	versionNumber := mux.Vars(r)["version"]
	applicationSlug := mux.Vars(r)["applicationSlug"]
	trackSlug := mux.Vars(r)["trackSlug"]

	version, err := models.GetVersionBySlugsAndNumber(applicationSlug, trackSlug, versionNumber)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	logo, contentType, err := storage.DownloadVersion(version.ApplicationId, version.TrackId, version.Id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	defer logo.Close()

	w.Header().Set("Content-Type", contentType)
	_, _ = io.Copy(w, logo)
}
