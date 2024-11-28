package userapp

import (
	"clanplan/server/bus/domain/userbus"
	"clanplan/server/foundation/web"
)

type Config struct {
	userbus *userbus.Business
}

func Routes(app *web.App){
	
}
