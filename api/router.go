package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func contentTypeJson() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, req)
		})
	}
}

func SetupApiRouter(router *mux.Router) {
	router.Methods("GET").Path("/api/admin/application").Handler(contentTypeJson()(http.HandlerFunc(getAllApplications)))
	router.Methods("POST").Path("/api/admin/application").Handler(contentTypeJson()(http.HandlerFunc(createApplication)))
	router.Methods("GET").Path("/api/admin/application/{id}").Handler(contentTypeJson()(http.HandlerFunc(getApplicationById)))
	router.Methods("PUT").Path("/api/admin/application/{id}").Handler(contentTypeJson()(http.HandlerFunc(updateApplication)))
	router.Methods("DELETE").Path("/api/admin/application/{id}").Handler(contentTypeJson()(http.HandlerFunc(deleteApplication)))
	router.Methods("POST").Path("/api/admin/application/{id}/logo").Handler(contentTypeJson()(http.HandlerFunc(deleteApplication)))

	router.Methods("GET").Path("/api/admin/application/{applicationId}/track").Handler(contentTypeJson()(http.HandlerFunc(getAllTracks)))
	router.Methods("POST").Path("/api/admin/application/{applicationId}/track").Handler(contentTypeJson()(http.HandlerFunc(createTrack)))
	router.Methods("GET").Path("/api/admin/application/{applicationId}/track/{id}").Handler(contentTypeJson()(http.HandlerFunc(getTrackById)))
	router.Methods("PUT").Path("/api/admin/application/{applicationId}/track/{id}").Handler(contentTypeJson()(http.HandlerFunc(updateTrack)))
	router.Methods("DELETE").Path("/api/admin/application/{applicationId}/track/{id}").Handler(contentTypeJson()(http.HandlerFunc(deleteTrack)))

	router.Methods("GET").Path("/api/admin/application/{applicationId}/track/{trackId}/version").Handler(contentTypeJson()(http.HandlerFunc(getAllVersions)))
	router.Methods("POST").Path("/api/admin/application/{applicationId}/track/{trackId}/version").Handler(contentTypeJson()(http.HandlerFunc(createVersion)))
	router.Methods("GET").Path("/api/admin/application/{applicationId}/track/{trackId}/version/{id}").Handler(contentTypeJson()(http.HandlerFunc(getVersionById)))
	router.Methods("DELETE").Path("/api/admin/application/{applicationId}/track/{trackId}/version/{id}").Handler(contentTypeJson()(http.HandlerFunc(deleteVersion)))
	router.Methods("POST").Path("/api/admin/application/{applicationId}/track/{trackId}/version/{id}/file").Handler(contentTypeJson()(http.HandlerFunc(uploadVersionBinary)))

	router.Methods("GET").Path("/api/admin/push-token").Handler(contentTypeJson()(http.HandlerFunc(getAllPushTokens)))
	router.Methods("POST").Path("/api/admin/push-token").Handler(contentTypeJson()(http.HandlerFunc(createPushToken)))
	router.Methods("GET").Path("/api/admin/push-token/{id}").Handler(contentTypeJson()(http.HandlerFunc(getPushTokenById)))
	router.Methods("PUT").Path("/api/admin/push-token/{id}").Handler(contentTypeJson()(http.HandlerFunc(updatePushToken)))
	router.Methods("DELETE").Path("/api/admin/push-token/{id}").Handler(contentTypeJson()(http.HandlerFunc(deletePushToken)))
}
