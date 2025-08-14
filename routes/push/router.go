package push

import (
	"jinya-releases/database"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func pushTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		if strings.HasPrefix(auth, "Bearer ") {
			token := strings.TrimPrefix(auth, "Bearer ")
			vars := mux.Vars(req)
			appSlug, exists := vars["applicationSlug"]
			pushTokenValid := database.CheckPushToken(token, appSlug)
			if exists && pushTokenValid {
				next.ServeHTTP(w, req)
				return
			}
		}

		w.WriteHeader(http.StatusUnauthorized)
	})
}

func SetupRouter(router *mux.Router) {
	router.
		Methods(http.MethodPost).
		Path("/api/push/{applicationSlug}/{trackSlug}/{versionNumber}").
		Handler(pushTokenMiddleware(http.HandlerFunc(pushVersion)))
}
