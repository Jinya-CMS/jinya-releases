package content

import (
	"fmt"
	"io"
	"jinya-releases/database"
	"jinya-releases/storage"
	"net/http"

	"github.com/gorilla/mux"
)

func getLogo(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]

	app, err := database.SelectOne[database.Application]("select * from application where slug = $1", slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	logo, contentType, contentLength, err := storage.DownloadLogo(app.Id.String())
	if err != nil {
		http.NotFound(w, r)
		return
	}

	defer logo.Close()

	w.Header().Set("Content-Length", fmt.Sprintf("%d", contentLength))
	w.Header().Set("Content-Type", contentType)
	_, _ = io.Copy(w, logo)
}
