package userbus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ardanlabs/service/business/sdk/order"
	"github.com/ardanlabs/service/business/sdk/page"
	"github.com/ardanlabs/service/foundation/logger"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound              = errors.New("user not found")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

type Storer interface {
	Count(ctx context.Context, filter QueryFilter) (int, error)
	Query(ctx context.Context, filter QueryFilter, order order.By, page page.Page) ([]User, error)
	QueryById(ctx context.Context, id uuid.UUID) (User, error)
	QueryOne(ctx context.Context, filter QueryFilter) (User, error)
	DeleteById(ctx context.Context, id string) error
	Delete(ctx context.Context, filter QueryFilter) error
	DeleteOne(ctx context.Context, filter QueryFilter) error
	Insert(ctx context.Context, u User) error
	Update(ctx context.Context, u User) error
}

type Business struct {
	storer Storer
	log    *logger.Logger
}

func NewBusiness(storer Storer, logger *logger.Logger) *Business {
	return &Business{storer: storer, log: logger}
}

func (b *Business) Create(ctx context.Context, nu NewUser) (User, error) {
	if nu.Password != nu.PasswordConfirm {
		passes := nu.Password + " " + nu.PasswordConfirm
		b.log.Error(ctx, "userbus.Create", passes, "Password and Password confirm must match")
		return User{}, fmt.Errorf("password and password confirm must match")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	now := time.Now()

	usr := User{
		Username:     nu.Username,
		Email:        nu.Email,
		PasswordHash: hash,
		Enabled:      true,
		DateCreated:  now,
		DateUpdated:  now,
	}

	if err := b.storer.Insert(ctx, usr); err != nil {
		b.log.Error(ctx, "userbus.Create", "Error inserting user", err)
		return User{}, fmt.Errorf("insert: %w", err)
	}
	return usr, nil
}

func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	count, err := b.storer.Count(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("count: %w", err)
	}
	return count, nil
}

func (b *Business) Query(ctx context.Context, filter QueryFilter, order order.By, page page.Page) ([]User, error) {
	usrs, err := b.storer.Query(ctx, filter, order, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return usrs, nil
}

func (b *Business) QueryById(ctx context.Context, id uuid.UUID) (User, error) {
	usr, err := b.storer.QueryById(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("querybyid: %w", err)
	}
	return usr, nil
}

func (b *Business) QueryOne(ctx context.Context, filter QueryFilter) (User, error) {
	usr, err := b.storer.QueryOne(ctx, filter)
	if err != nil {
		return User{}, fmt.Errorf("queryone: %w", err)
	}
	return usr, nil
}

func (b *Business) Update(ctx context.Context, usr User, uu UpdateUser) (User, error) {
	if uu.Password != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, fmt.Errorf("generatefrompassword: %w", err)
		}
		usr.PasswordHash = pw
	}

	if uu.Email != nil {
		usr.Email = *uu.Email
	}

	if uu.Enabled != nil {
		usr.Enabled = *uu.Enabled
	}

	if uu.Username != nil {
		usr.Username = *uu.Username
	}

	if uu.Email != nil {
		usr.Email = *uu.Email
	}

	usr.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, usr); err != nil {
		return User{}, fmt.Errorf("update: %w", err)
	}

	return usr, nil
}

func (b *Business) DeleteById(ctx context.Context, userID uuid.UUID) error {
	if err := b.storer.DeleteById(ctx, userID.String()); err != nil {
		return fmt.Errorf("deletebyid: %w", err)
	}
	return nil
}

func (b *Business) Delete(ctx context.Context, filter QueryFilter) error {
	if err := b.storer.Delete(ctx, filter); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (b *Business) DeleteOne(ctx context.Context, filter QueryFilter) error {
	if err := b.storer.DeleteOne(ctx, filter); err != nil {
		return fmt.Errorf("deleteone: %w", err)
	}
	return nil
}
