package api

import (
	"context"
	"jinya-releases/config"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

func contentTypeJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, req)
	})
}

func SetupRouter(router *mux.Router) {
	ctx := context.Background()

	zitadelConfig := oauth.WithIntrospection[*oauth.IntrospectionContext](oauth.ClientIDSecretIntrospectionAuthentication(config.LoadedConfiguration.OidcServerClientId, config.LoadedConfiguration.OidcServerClientSecret))
	authZ, err := authorization.New(ctx, zitadel.New(config.LoadedConfiguration.OidcDomain), zitadelConfig)

	if err != nil {
		panic(err)
	}

	mw := middleware.New(authZ)

	adminRouter := router.PathPrefix("/api/admin").Subrouter()
	adminRouter.
		Methods(http.MethodGet).
		Path("/application").
		HandlerFunc(getAllApplications)
	adminRouter.
		Methods(http.MethodPost).
		Path("/application").
		HandlerFunc(createApplication)
	adminRouter.
		Methods(http.MethodGet).
		Path("/application/{id}").
		HandlerFunc(getApplicationById)
	adminRouter.
		Methods(http.MethodPut).
		Path("/application/{id}").
		HandlerFunc(updateApplication)
	adminRouter.
		Methods(http.MethodDelete).
		Path("/application/{id}").
		HandlerFunc(deleteApplication)
	adminRouter.
		Methods(http.MethodPost).
		Path("/application/{id}/logo").
		HandlerFunc(uploadLogo)
	adminRouter.
		Methods(http.MethodDelete).
		Path("/application/{id}/token").
		HandlerFunc(resetToken)
	adminRouter.
		Methods(http.MethodPost).
		Path("/application/{id}/token").
		HandlerFunc(createToken)

	adminRouter.
		Methods(http.MethodGet).
		Path("/application/{applicationId}/track").
		HandlerFunc(getAllTracks)
	adminRouter.
		Methods(http.MethodPost).
		Path("/application/{applicationId}/track").
		HandlerFunc(createTrack)
	adminRouter.
		Methods(http.MethodGet).
		Path("/application/{applicationId}/track/{id}").
		HandlerFunc(getTrackById)
	adminRouter.
		Methods(http.MethodPut).
		Path("/application/{applicationId}/track/{id}").
		HandlerFunc(updateTrack)
	adminRouter.
		Methods(http.MethodDelete).
		Path("/application/{applicationId}/track/{id}").
		HandlerFunc(deleteTrack)

	adminRouter.
		Methods(http.MethodGet).
		Path("/application/{applicationId}/track/{trackId}/version").
		HandlerFunc(getAllVersions)
	adminRouter.
		Methods(http.MethodGet).
		Path("/application/{applicationId}/track/{trackId}/version/{id}").
		HandlerFunc(getVersionById)
	adminRouter.
		Methods(http.MethodDelete).
		Path("/application/{applicationId}/track/{trackId}/version/{id}").
		HandlerFunc(deleteVersion)
	adminRouter.
		Methods(http.MethodPost).
		Path("/application/{applicationId}/track/{trackId}/version/{version}").
		HandlerFunc(uploadVersion)

	adminRouter.Use(mw.RequireAuthorization(), contentTypeJson)
}
