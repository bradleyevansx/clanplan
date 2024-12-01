package userapp

import "clanplan/server/bus/domain/userbus"

var orderByFields = map[string]string{
	"user_id": userbus.OrderByID,
	"name":    userbus.OrderByName,
	"email":   userbus.OrderByEmail,
	"enabled": userbus.OrderByEnabled,
}
