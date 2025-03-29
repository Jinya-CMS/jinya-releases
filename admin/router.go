package admin

import (
	"github.com/gorilla/mux"
	"jinya-releases/admin/web"
)

func SetupRouter(router *mux.Router) {
	web.SetupRouter(router)
}
