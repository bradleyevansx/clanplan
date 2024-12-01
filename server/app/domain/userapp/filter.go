package userapp

import (
	"clanplan/server/bus/domain/userbus"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/ardanlabs/service/app/sdk/errs"
	"github.com/ardanlabs/service/business/types/name"
	"github.com/google/uuid"
)

func parseQueryParams(r *http.Request) (queryParams, error) {
	values := r.URL.Query()

	filter := queryParams{
		Page:             values.Get("page"),
		Row:              values.Get("row"),
		Order:            values.Get("order"),
		ID:               values.Get("user_id"),
		Username:         values.Get("username"),
		Email:            values.Get("email"),
		StartCreatedDate: values.Get("start_date_created"),
		EndCreatedDate:   values.Get("end_date_created"),
	}

	fmt.Println(filter)

	return filter, nil
}

func parseFilter(qp queryParams) (userbus.QueryFilter, error) {
	var filter userbus.QueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			return userbus.QueryFilter{}, errs.NewFieldsError("user_id", err)
		}
		filter.ID = &id
	}

	if qp.Username != "" {
		name, err := name.Parse(qp.Username)
		if err != nil {
			return userbus.QueryFilter{}, errs.NewFieldsError("username", err)
		}
		filter.Username = &name
	}

	if qp.Email != "" {
		addr, err := mail.ParseAddress(qp.Email)
		if err != nil {
			return userbus.QueryFilter{}, errs.NewFieldsError("email", err)
		}
		filter.Email = addr
	}

	if qp.StartCreatedDate != "" {
		t, err := time.Parse(time.RFC3339, qp.StartCreatedDate)
		if err != nil {
			return userbus.QueryFilter{}, errs.NewFieldsError("start_date_created", err)
		}
		filter.StartCreatedDate = &t
	}

	if qp.EndCreatedDate != "" {
		t, err := time.Parse(time.RFC3339, qp.EndCreatedDate)
		if err != nil {
			return userbus.QueryFilter{}, errs.NewFieldsError("end_date_created", err)
		}
		filter.EndCreatedDate = &t
	}

	return filter, nil
}
