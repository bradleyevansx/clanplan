package userbus

import (
	"context"
	"fmt"
	"time"

	"github.com/ardanlabs/service/foundation/logger"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Storer interface {
	Query(ctx context.Context, filter QueryFilter) ([]User, error)
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
func (b *Business) Query(ctx context.Context, filter QueryFilter) ([]User, error) {
	usrs, err := b.storer.Query(ctx, filter)
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

func (b *Business) Update(ctx context.Context, usr User, uu UpdateUser) error {

}
