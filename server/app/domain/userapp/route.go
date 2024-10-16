package userapp

import (
	"clanplan/server/bus/domain/userbus"
	web "clanplan/server/foundation"
)

type Config struct {
	userbus *userbus.Business
}

func Routes(app *web.App){
	
}