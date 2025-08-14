package routes

import (
	"jinya-releases/routes/admin"
	"jinya-releases/routes/content"
	"jinya-releases/routes/frontend"
	"jinya-releases/routes/push"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router) {
	admin.SetupRouter(router)
	push.SetupRouter(router)
	content.SetupRouter(router)
	frontend.SetupRouter(router)
}
