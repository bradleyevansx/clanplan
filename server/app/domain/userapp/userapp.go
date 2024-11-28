package userapp

import (
	"clanplan/server/app/sdk/errs"
	"clanplan/server/app/sdk/errs/query"
	"clanplan/server/bus/domain/userbus"
	"clanplan/server/foundation/web"
	"context"
	"net/http"
)

type app struct {
	userBus *userbus.Business
}

func newApp(userBus *userbus.Business) *app {
	return &app{
		userBus: userBus,
	}
}

//type Storer interface {
//	Query(ctx context.Context, filter QueryFilter) ([]User, error)
//	QueryById(ctx context.Context, id uuid.UUID) (User, error)
//	QueryOne(ctx context.Context, filter QueryFilter) (User, error)
//	DeleteById(ctx context.Context, id string) error
//	Delete(ctx context.Context, filter QueryFilter) error
//	DeleteOne(ctx context.Context, filter QueryFilter) error
//	Insert(ctx context.Context, u User) error
//	Update(ctx context.Context, u User) error
//}

func (a *app) query(ctx context.Context, r *http.Request) web.Encoder {
	qr, err := parseQueryParams(r)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	filter, err := parseFilter(qr)
	if err != nil {
		return err.(errs.FieldErrors)
	}

	users, err := a.userBus.Query(ctx, filter)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	return query.NewResult(toAppUsers(users))
}

func (a *app) delete(ctx context.Context, r *http.Request) web.Encoder {
	params, err := parseQueryParams(r)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	filter, err := parseFilter(params)
	if err != nil {
		return err.(errs.FieldErrors)
	}

	err = a.userBus.Delete(ctx, filter)
	if err != nil {
		return errs.Newf(errs.Internal, "delete: %s", err)
	}
	return nil
}

func (a *app) insert(ctx context.Context, r *http.Request) web.Encoder {
	var app NewUser

	err := web.Decode(r, &app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nu, err := toBus(app)
}
