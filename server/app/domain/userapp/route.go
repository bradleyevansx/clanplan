package userapp

import (
	"clanplan/server/bus/domain/userbus"
	"clanplan/server/foundation/web"
)

type Config struct {
	Userbus *userbus.Business
}

func Routes(app *web.App, cfg Config) {
	group := "v1"

	api := newApp(cfg.Userbus)

	app.HandlerFunc("GET", group, "/users/", api.query)
	app.HandlerFunc("DELETE", group, "/users/", api.delete)
}
