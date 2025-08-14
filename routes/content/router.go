package content

import "github.com/gorilla/mux"

func SetupContentRouter(router *mux.Router) {
	router.Methods("GET").Path("/content/logo/{slug}").HandlerFunc(GetLogo)
	router.Methods("GET").Path("/content/version/{applicationSlug}/{trackSlug}/{version}").HandlerFunc(GetVersion)
}
