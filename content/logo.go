package content

import (
	"fmt"
	"io"
	"jinya-releases/database/models"
	"jinya-releases/storage"
	"net/http"

	"github.com/gorilla/mux"
)

func GetLogo(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]

	app, err := models.GetApplicationBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	logo, contentType, contentLength, err := storage.DownloadLogo(app.Id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	defer logo.Close()

	w.Header().Set("Content-Length", fmt.Sprintf("%d", contentLength))
	w.Header().Set("Content-Type", contentType)
	_, _ = io.Copy(w, logo)
}
