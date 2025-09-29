package frontend

import (
	"embed"
	"html/template"
	"jinya-releases/config"
	"jinya-releases/database"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

//go:embed templates
var templates embed.FS

func render(w http.ResponseWriter, name string, data any) {
	tmpl, err := template.New("layout").Funcs(map[string]any{
		"toTimeString": func(time time.Time) string {
			return time.Format("2006-01-02")
		},
	}).ParseFS(templates, "templates/layout.gohtml", "templates/"+name+".gohtml")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getHomePage(w http.ResponseWriter, r *http.Request) {
	apps, err := database.Select[database.Application]("select * from application order by name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render(w, "home", struct {
		Applications []database.Application
	}{
		Applications: apps,
	})
}

func getImprintPage(w http.ResponseWriter, r *http.Request) {
	render(w, "imprint", nil)
}

func getDataProtectionPage(w http.ResponseWriter, r *http.Request) {
	render(w, "data-protection", nil)
}

func getApplicationPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	applicationSlug := vars["applicationSlug"]

	app, err := database.SelectOne[database.Application]("select * from application where slug = $1", applicationSlug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tracks, err := database.Select[database.Track]("select * from track where application_id = $1 order by name", app.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(tracks) == 1 {
		http.Redirect(w, r, "/"+applicationSlug+"/"+tracks[0].Slug, http.StatusSeeOther)
		return
	}

	tracksWithVersions := make([]database.Track, 0)
	for _, track := range tracks {
		versionsCount, err := database.GetDbMap().SelectInt("select count(*) from version where track_id = $1", track.Id)
		if err != nil || versionsCount == 0 {
			continue
		}

		tracksWithVersions = append(tracksWithVersions, track)
	}

	slices.SortFunc(tracksWithVersions, func(a, b database.Track) int {
		if a.IsDefault {
			return -1
		}

		if b.IsDefault {
			return 1
		}

		return strings.Compare(a.Name, b.Name)
	})

	render(w, "application", struct {
		Application database.Application
		Tracks      []database.Track
	}{
		Application: app,
		Tracks:      tracksWithVersions,
	})
}

func getTrackPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	applicationSlug := vars["applicationSlug"]
	trackSlug := vars["trackSlug"]

	app, err := database.SelectOne[database.Application]("select * from application where slug = $1", applicationSlug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	track, err := database.SelectOne[database.Track]("select * from track where application_id = $1 and slug = $2", app.Id, trackSlug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	versions, err := database.Select[database.Version](`
select v.*, $1 || '/content/version/' || a.slug || '/' || t.slug || '/' || v.version as url
from version v
		inner join application a on a.id = v.application_id
		inner join track t on v.track_id = t.id
where a.slug = $2 and t.slug = $3
order by upload_date desc
		`, config.LoadedConfiguration.ServerUrl, applicationSlug, trackSlug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render(w, "track", struct {
		Application database.Application
		Track       database.Track
		Versions    []database.Version
	}{
		Application: app,
		Track:       track,
		Versions:    versions,
	})
}
