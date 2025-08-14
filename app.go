package main

import (
	"embed"
	"jinya-releases/config"
	migrator "jinya-releases/database/migrations"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var (
	//go:embed openapi/admin
	adminOpenapi embed.FS
	//go:embed static
	static embed.FS
)

type SpaHandler struct {
	embedFS   embed.FS
	indexPath string
}

func (handler SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fullPath := strings.TrimPrefix(r.URL.Path, "/")
	file, err := handler.embedFS.Open(fullPath)
	if err != nil {
		http.ServeFileFS(w, r, handler.embedFS, handler.indexPath)
		return
	}

	if fi, err := file.Stat(); err != nil || fi.IsDir() {
		http.ServeFileFS(w, r, handler.embedFS, handler.indexPath)
		return
	}

	http.ServeFileFS(w, r, handler.embedFS, fullPath)
}

func main() {
	err := config.LoadConfiguration()
	if err != nil {
		panic(err)
	}

	err = migrator.Migrate()
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	router.PathPrefix("/openapi/admin").Handler(SpaHandler{
		embedFS:   adminOpenapi,
		indexPath: "openapi/admin/index.html",
	})

	router.PathPrefix("/static/").Handler(http.FileServerFS(static))

	log.Println("Serving at localhost:8090...")
	err = http.ListenAndServe(":8090", router)
	if err != nil {
		panic(err)
	}
}
