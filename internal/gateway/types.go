package gateway

import (
	"context"
)

//go:generate mockgen -source=types.go -destination=mocks_test.go -package=gateway

type authOperations interface {
	Registration(ctx context.Context, email string) error
	ConfirmationOTP(ctx context.Context, code string, email string) error
	Confirmation(ctx context.Context, email string, password string, login string) error
	Login(ctx context.Context, login string, password string) (string, error)
}
