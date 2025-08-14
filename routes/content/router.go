package content

import "github.com/gorilla/mux"

func SetupRouter(router *mux.Router) {
	router.Methods("GET").Path("/content/logo/{slug}").HandlerFunc(getLogo)
	router.Methods("GET").Path("/content/version/{applicationSlug}/{trackSlug}/{version}").HandlerFunc(getVersion)
}
