package content

import (
	"github.com/gorilla/mux"
	"io"
	"jinya-releases/database/models"
	"jinya-releases/storage"
	"net/http"
)

func GetLogo(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]

	app, err := models.GetApplicationBySlug(slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	logo, err := storage.DownloadLogo(app.Id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer logo.Close()

	_, _ = io.Copy(w, logo)
}
