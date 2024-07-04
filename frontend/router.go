package frontend

import (
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func checkForPageOrJson(expectPage bool) func(request *http.Request, match *mux.RouteMatch) bool {
	return func(request *http.Request, match *mux.RouteMatch) bool {
		if expectPage {
			return strings.Contains(request.Header.Get("Accept"), "text/html")
		}

		return true
	}
}

func SetupFrontendRouter(router *mux.Router) {
	router.Methods("GET").Path("/").HandlerFunc(getHomePage)
	router.Methods("GET").Path("/imprint").HandlerFunc(getImprintPage)
	router.Methods("GET").Path("/data-protection").HandlerFunc(getDataProtectionPage)

	router.Methods("GET").Path("/{applicationSlug}").MatcherFunc(checkForPageOrJson(true)).HandlerFunc(getApplicationPage)
	router.Methods("GET").Path("/{applicationSlug}/{trackSlug}").MatcherFunc(checkForPageOrJson(true)).HandlerFunc(getTrackPage)

	router.Methods("GET").Path("/{applicationSlug}").MatcherFunc(checkForPageOrJson(false)).HandlerFunc(getApplicationJson)
	router.Methods("GET").Path("/{applicationSlug}/{trackSlug}").MatcherFunc(checkForPageOrJson(false)).HandlerFunc(getTrackJson)
}
