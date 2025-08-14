package content

import (
	"fmt"
	"io"
	"jinya-releases/database"
	"jinya-releases/storage"
	"net/http"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gorilla/mux"
)

func getVersion(w http.ResponseWriter, r *http.Request) {
	versionNumber := mux.Vars(r)["version"]
	applicationSlug := mux.Vars(r)["applicationSlug"]
	trackSlug := mux.Vars(r)["trackSlug"]

	app, err := database.SelectOne[database.Application]("select * from application where slug = $1", applicationSlug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	version, err := database.SelectOne[database.Version](`
select v.*
from version v
         inner join application a on a.id = v.application_id
         inner join track t on t.id = v.track_id
where v.version = $1
  and a.slug = $2
  and t.slug = $3
`, versionNumber, applicationSlug, trackSlug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	v, contentType, contentLength, err := storage.DownloadVersion(version.ApplicationId, version.TrackId, version.Id.String())
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
