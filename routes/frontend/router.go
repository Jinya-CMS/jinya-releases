package frontend

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func checkForPageOrJson(expectPage bool) func(request *http.Request, match *mux.RouteMatch) bool {
	return func(request *http.Request, match *mux.RouteMatch) bool {
		if expectPage {
			return strings.Contains(request.Header.Get("Accept"), "text/html")
		}

		return true
	}
}

func SetupRouter(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/").HandlerFunc(getHomePage)
	router.Methods(http.MethodGet).Path("/imprint").HandlerFunc(getImprintPage)
	router.Methods(http.MethodGet).Path("/data-protection").HandlerFunc(getDataProtectionPage)

	router.
		Methods(http.MethodGet).
		Headers("Accept", "text/html").
		Path("/{applicationSlug}").
		HandlerFunc(getApplicationPage)
	router.
		Methods(http.MethodGet).
		Headers("Accept", "text/html").
		Path("/{applicationSlug}/{trackSlug}").
		HandlerFunc(getTrackPage)

	router.
		Methods(http.MethodGet).
		Path("/{applicationSlug}").
		HandlerFunc(getApplicationJson)
	router.
		Methods(http.MethodGet).
		Path("/{applicationSlug}/{trackSlug}").
		HandlerFunc(getTrackJson)
}
