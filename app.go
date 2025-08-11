package main

import (
	"embed"
	"html/template"
	"jinya-releases/admin"
	"jinya-releases/api"
	"jinya-releases/config"
	"jinya-releases/content"
	migrator "jinya-releases/database/migrations"
	"jinya-releases/frontend"
	"jinya-releases/importer"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/text/language"
)

var (
	//go:embed openapi/admin
	adminOpenapi embed.FS
	//go:embed openapi
	openapi embed.FS
	//go:embed static
	static embed.FS
)

type SpaHandler struct {
	embedFS      embed.FS
	indexPath    string
	fsPrefixPath string
	templated    bool
	templateData any
}

type LanguageHandler struct {
	defaultLang     language.Tag
	langPathMapping map[language.Tag]string
}

func (handler SpaHandler) serveTemplated(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFS(handler.embedFS, handler.indexPath)
	if err != nil {
		http.Error(w, "Failed to get admin page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, handler.templateData)
	if err != nil {
		http.Error(w, "Failed to get admin page", http.StatusInternalServerError)
		return
	}
}

func (handler SpaHandler) servePlain(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, handler.embedFS, handler.indexPath)
}

func (handler SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fullPath := strings.TrimPrefix(path.Join(handler.fsPrefixPath, r.URL.Path), "/")
	file, err := handler.embedFS.Open(fullPath)
	if err != nil {
		if handler.templated {
			handler.serveTemplated(w, r)
		} else {
			handler.servePlain(w, r)
		}
		return
	}

	if fi, err := file.Stat(); err != nil || fi.IsDir() {
		if handler.templated {
			handler.serveTemplated(w, r)
		} else {
			handler.servePlain(w, r)
		}
		return
	}

	http.ServeFileFS(w, r, handler.embedFS, fullPath)
}

func (handler LanguageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	oldPath := "/" + r.URL.Path + "?" + r.URL.RawQuery
	acceptLanguage, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	if err != nil {
		http.Redirect(w, r, handler.langPathMapping[handler.defaultLang]+oldPath, http.StatusFound)
		return
	}

	localMap := map[string]string{}
	for tag, p := range handler.langPathMapping {
		b, _ := tag.Base()
		localMap[b.ISO3()] = p
	}

	for _, lang := range acceptLanguage {
		b, _ := lang.Base()
		if p, exists := localMap[b.ISO3()]; exists {
			http.Redirect(w, r, p+oldPath, http.StatusFound)
			return
		}
	}

	http.Redirect(w, r, handler.langPathMapping[handler.defaultLang]+oldPath, http.StatusFound)
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

	args := os.Args
	if len(args) > 1 {
		if args[1] == "import" {
			importer.ImportFromFileSystem(args[2:])
			return
		} else {
			log.Fatalf("Unknown argument '%s'", args[1])
		}
	} else {
		router := mux.NewRouter()

		router.PathPrefix("/openapi/admin").Handler(SpaHandler{
			embedFS:      adminOpenapi,
			indexPath:    "openapi/admin/index.html",
			fsPrefixPath: "",
			templated:    false,
		})
		//router.PathPrefix("/admin/de").Handler(http.StripPrefix("/admin/de", SpaHandler{
		//	embedFS:      angularAdmin,
		//	indexPath:    "angular/dist/admin/browser/de/index.html",
		//	fsPrefixPath: "angular/dist/admin/browser/de",
		//	templated:    true,
		//	templateData: config.LoadedConfiguration,
		//}))
		//router.PathPrefix("/admin/en").Handler(http.StripPrefix("/admin/en", SpaHandler{
		//	embedFS:      angularAdmin,
		//	indexPath:    "angular/dist/admin/browser/en/index.html",
		//	fsPrefixPath: "angular/dist/admin/browser/en",
		//	templated:    true,
		//	templateData: config.LoadedConfiguration,
		//}))
		//router.PathPrefix("/admin").Handler(http.StripPrefix("/admin", LanguageHandler{
		//	defaultLang: language.English,
		//	langPathMapping: map[language.Tag]string{
		//		language.English: "/admin/en",
		//		language.German:  "/admin/de",
		//	},
		//}))

		router.PathPrefix("/static/").Handler(http.FileServerFS(static))

		admin.SetupRouter(router)
		api.SetupApiRouter(router)
		content.SetupContentRouter(router)
		frontend.SetupFrontendRouter(router)

		log.Println("Serving at localhost:8090...")
		err = http.ListenAndServe(":8090", router)
		if err != nil {
			panic(err)
		}
	}
}
