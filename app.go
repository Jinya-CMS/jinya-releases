package main

import (
	"embed"
	"github.com/gorilla/mux"
	"jinya-releases/api"
	"jinya-releases/config"
	migrator "jinya-releases/database/migrations"
	dbMigrations "jinya-releases/migrations"
	"log"
	"net/http"
	"path"
	"strings"
)

var (
	//go:embed openapi/admin
	adminOpenapi embed.FS
	//go:embed openapi
	openapi embed.FS
	//go:embed static
	static embed.FS
)

var migrations = []migrator.Migration{
	dbMigrations.CreateApplicationTable{},
}

type SpaHandler struct {
	embedFS      embed.FS
	indexPath    string
	fsPrefixPath string
}

func (handler SpaHandler) servePlain(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, handler.embedFS, handler.indexPath)
}

func (handler SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fullPath := strings.TrimPrefix(path.Join(handler.fsPrefixPath, r.URL.Path), "/")
	file, err := handler.embedFS.Open(fullPath)
	if err != nil {
		handler.servePlain(w, r)
		return
	}

	if fi, err := file.Stat(); err != nil || fi.IsDir() {
		handler.servePlain(w, r)
		return
	}

	http.ServeFileFS(w, r, handler.embedFS, fullPath)
}

func main() {
	err := config.LoadConfiguration()
	if err != nil {
		panic(err)
	}

	err = migrator.Migrate(migrations)
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	router.PathPrefix("/openapi/admin").Handler(SpaHandler{
		embedFS:      adminOpenapi,
		indexPath:    "openapi/admin/index.html",
		fsPrefixPath: "",
	})
	router.PathPrefix("/openapi").Handler(SpaHandler{
		embedFS:      openapi,
		indexPath:    "openapi/index.html",
		fsPrefixPath: "",
	})

	router.PathPrefix("/static/").Handler(http.FileServerFS(static))

	api.SetupApiRouter(router)

	log.Println("Serving at localhost:8090...")
	err = http.ListenAndServe(":8090", router)
	if err != nil {
		panic(err)
	}
}
