package frontend

import (
	"encoding/json"
	"jinya-releases/database"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iancoleman/orderedmap"
)

func versionsToMap(versions []database.Version) *orderedmap.OrderedMap {
	versionMap := orderedmap.New()
	for _, version := range versions {
		versionMap.Set(version.Version, version.Url)
	}

	return versionMap
}

func getApplicationJson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	applicationSlug := vars["applicationSlug"]
	encoder := json.NewEncoder(w)

	versions, err := database.Select[database.Version](`
select v.*
from version v
         inner join application a on a.id = v.application_id
         inner join track t on v.track_id = t.id
where a.slug = $1 and t.is_default
`, applicationSlug)
	if err != nil {
		versions = []database.Version{}
	}

	_ = encoder.Encode(versionsToMap(versions))
}

func getTrackJson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	applicationSlug := vars["applicationSlug"]
	trackSlug := vars["trackSlug"]

	encoder := json.NewEncoder(w)

	versions, err := database.Select[database.Version](`
select v.*
from version v
         inner join application a on a.id = v.application_id
         inner join track t on v.track_id = t.id
where a.slug = $1 and t.slug = $2 and t.is_default
`, applicationSlug, trackSlug)
	if err != nil {
		versions = []database.Version{}
	}

	_ = encoder.Encode(versionsToMap(versions))
}
