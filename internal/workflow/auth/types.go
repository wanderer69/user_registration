package auth

import (
	"context"

	"github.com/wanderer69/user_registration/internal/entity"
)

//go:generate mockgen -source=types.go -destination=mocks.go -package=auth

type userRepository interface {
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	ConfirmationUpdate(ctx context.Context, user *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByLogin(ctx context.Context, login string) (*entity.User, error)
	GetByUUID(ctx context.Context, uuid string) (*entity.User, error)
	GetByRegistrationCode(ctx context.Context, code string) (*entity.User, error)
	DeleteByUUID(ctx context.Context, userUUID string) error
}

type mailService interface {
	Send(email string, subject string, message string, fromName string) error
}
