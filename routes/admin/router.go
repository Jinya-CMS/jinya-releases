package admin

import (
	"jinya-releases/routes/admin/api"
	"jinya-releases/routes/admin/web"

	"github.com/gorilla/mux"
)

func SetupRouter(router *mux.Router) {
	api.SetupRouter(router)
	web.SetupRouter(router)
}
