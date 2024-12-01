package userapp

import (
	"clanplan/server/bus/domain/userbus"
	"clanplan/server/foundation/web"
	"context"
	"errors"
	"net/http"

	"github.com/ardanlabs/service/app/sdk/errs"
	"github.com/ardanlabs/service/app/sdk/query"
	"github.com/ardanlabs/service/business/sdk/order"
	"github.com/ardanlabs/service/business/sdk/page"
)

type app struct {
	userBus *userbus.Business
}

func newApp(userBus *userbus.Business) *app {
	return &app{
		userBus: userBus,
	}
}

func (a *app) query(ctx context.Context, r *http.Request) web.Encoder {
	qr, err := parseQueryParams(r)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	page, err := page.Parse(qr.Page, qr.Row)
	if err != nil {
		return errs.NewFieldsError("page", err)
	}

	filter, err := parseFilter(qr)
	if err != nil {
		return err.(errs.FieldErrors)
	}

	order, err := order.Parse(orderByFields, qr.Order, userbus.DefaultOrderBy)
	if err != nil {
		return errs.NewFieldsError("order", err)
	}

	users, err := a.userBus.Query(ctx, filter, order, page)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	count, err := a.userBus.Count(ctx, filter)
	if err != nil {
		return errs.Newf(errs.Internal, "count: %s", err)
	}

	return query.NewResult(toAppUsers(users), count, page)
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

	nu, err := toNewBusUser(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	usr, err := a.userBus.Create(ctx, nu)
	if err != nil {
		if errors.Is(err, userbus.ErrUniqueEmail) {
			return errs.New(errs.Aborted, userbus.ErrUniqueEmail)
		}
		return errs.Newf(errs.Internal, "create: usr[%+v]: %s", usr, err)
	}

	return toAppUser(usr)
}
