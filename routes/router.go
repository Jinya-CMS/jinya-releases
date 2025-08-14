package routes

import (
	"jinya-releases/routes/admin"
	"jinya-releases/routes/api"
	"jinya-releases/routes/content"
	"jinya-releases/routes/frontend"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router) {
	admin.SetupRouter(router)
	api.SetupApiRouter(router)
	content.SetupContentRouter(router)
	frontend.SetupFrontendRouter(router)
}
