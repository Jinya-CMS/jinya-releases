package content

import (
	"fmt"
	"io"
	"jinya-releases/database/models"
	"jinya-releases/storage"
	"net/http"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gorilla/mux"
)

func GetVersion(w http.ResponseWriter, r *http.Request) {
	versionNumber := mux.Vars(r)["version"]
	applicationSlug := mux.Vars(r)["applicationSlug"]
	trackSlug := mux.Vars(r)["trackSlug"]

	app, err := models.GetApplicationBySlug(applicationSlug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	version, err := models.GetVersionBySlugsAndNumber(applicationSlug, trackSlug, versionNumber)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	v, contentType, contentLength, err := storage.DownloadVersion(version.ApplicationId, version.TrackId, version.Id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	defer v.Close()

	mime := mimetype.Lookup(contentType)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", contentLength))
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s %s.%s"`, app.Name, version.Version, mime.Extension()))
	_, _ = io.Copy(w, v)
}
