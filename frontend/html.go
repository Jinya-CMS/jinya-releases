package frontend

import (
	"embed"
	"github.com/gorilla/mux"
	"html/template"
	"jinya-releases/database/models"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"
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
	apps, err := models.GetAllApplications()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render(w, "home", struct {
		Applications []models.Application
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

	if len(tracks) == 1 {
		http.Redirect(w, r, "/"+applicationSlug+"/"+tracks[0].Slug, http.StatusSeeOther)
		return
	}

	tracksWithVersions := make([]models.Track, 0)
	for _, track := range tracks {
		versions, err := models.GetAllVersions(track.ApplicationId, track.Id)
		if err != nil || len(versions) == 0 {
			continue
		}

		tracksWithVersions = append(tracksWithVersions, track)
	}

	slices.SortFunc(tracksWithVersions, func(a, b models.Track) int {
		if a.IsDefault {
			return 1
		}

		if b.IsDefault {
			return -1
		}

		return strings.Compare(a.Name, b.Name)
	})

	render(w, "application", struct {
		Application *models.Application
		Tracks      []models.Track
	}{
		Application: app,
		Tracks:      tracks,
	})
}

func getTrackPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	applicationSlug := vars["applicationSlug"]
	trackSlug := vars["trackSlug"]

	app, err := models.GetApplicationBySlug(applicationSlug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	track, err := models.GetTrackBySlug(trackSlug, app.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	versions, err := models.GetAllVersions(app.Id, track.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render(w, "track", struct {
		Application *models.Application
		Track       *models.Track
		Versions    []models.Version
	}{
		Application: app,
		Track:       track,
		Versions:    versions,
	})
}
