package admin

import (
	"jinya-releases/admin/web"

	"github.com/gorilla/mux"
)

func SetupRouter(router *mux.Router) {
	web.SetupRouter(router)
}
