package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/gorilla/mux"
	"github.com/zitadel/oidc/v3/pkg/client"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
	"jinya-releases/config"
	"jinya-releases/database/models"
	"net/http"
	"strings"
)

func contentTypeJson() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, req)
		})
	}
}

func pushTokenMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			auth := req.Header.Get("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				token := strings.TrimPrefix(auth, "Bearer ")
				vars := mux.Vars(req)
				appSlug, exists := vars["applicationSlug"]
				if exists && models.CheckPushToken(token, appSlug) {
					next.ServeHTTP(w, req)
					return
				}
			}

			w.WriteHeader(http.StatusUnauthorized)
		})
	}
}

func SetupApiRouter(router *mux.Router) {
	ctx := context.Background()
	encryptionKey := make([]byte, 32)
	_, err := rand.Read(encryptionKey)
	if err != nil {
		panic(err)
	}

	keyFileData, err := base64.StdEncoding.DecodeString(config.LoadedConfiguration.OpenIDKeyFileData)
	if err != nil {
		panic(err)
	}

	keyFile, err := client.ConfigFromKeyFileData(keyFileData)
	zitadelConfig := oauth.WithIntrospection[*oauth.IntrospectionContext](oauth.JWTProfileIntrospectionAuthentication(keyFile))
	authZ, err := authorization.New(ctx, zitadel.New(config.LoadedConfiguration.OpenIDDomain), zitadelConfig)

	if err != nil {
		panic(err)
	}

	mw := middleware.New(authZ)

	router.Methods("GET").Path("/api/admin/application").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(getAllApplications))))
	router.Methods("POST").Path("/api/admin/application").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(createApplication))))
	router.Methods("GET").Path("/api/admin/application/{id}").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(getApplicationById))))
	router.Methods("PUT").Path("/api/admin/application/{id}").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(updateApplication))))
	router.Methods("DELETE").Path("/api/admin/application/{id}").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(deleteApplication))))
	router.Methods("POST").Path("/api/admin/application/{id}/logo").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(uploadLogo))))
	router.Methods("DELETE").Path("/api/admin/application/{id}/token").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(resetToken))))
	router.Methods("POST").Path("/api/admin/application/{id}/token").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(createToken))))

	router.Methods("GET").Path("/api/admin/application/{applicationId}/track").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(getAllTracks))))
	router.Methods("POST").Path("/api/admin/application/{applicationId}/track").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(createTrack))))
	router.Methods("GET").Path("/api/admin/application/{applicationId}/track/{id}").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(getTrackById))))
	router.Methods("PUT").Path("/api/admin/application/{applicationId}/track/{id}").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(updateTrack))))
	router.Methods("DELETE").Path("/api/admin/application/{applicationId}/track/{id}").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(deleteTrack))))

	router.Methods("GET").Path("/api/admin/application/{applicationId}/track/{trackId}/version").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(getAllVersions))))
	router.Methods("GET").Path("/api/admin/application/{applicationId}/track/{trackId}/version/{id}").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(getVersionById))))
	router.Methods("DELETE").Path("/api/admin/application/{applicationId}/track/{trackId}/version/{id}").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(deleteVersion))))
	router.Methods("POST").Path("/api/admin/application/{applicationId}/track/{trackId}/version/{versionNumber}").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(uploadVersion))))

	router.Methods("POST").Path("/api/push/{applicationSlug}/{trackSlug}/{versionNumber}").Handler(pushTokenMiddleware()(http.HandlerFunc(pushVersion)))
}
