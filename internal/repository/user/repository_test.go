package user

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/wanderer69/user_registration/internal/entity"
	"github.com/wanderer69/user_registration/pkg/tools/tests"
)

func TestRepository(t *testing.T) {
	ctx := context.Background()
	dao, err := tests.InitDAO("../../../migrations")
	require.NoError(t, err)

	userRepo := NewRepository(dao)
	user1 := &entity.User{
		UUID:  uuid.NewString(),
		Login: "user1",
		Email: "user1@example.com",
	}
	require.NoError(t, userRepo.Create(ctx, user1))
	require.ErrorContains(t, userRepo.Create(ctx, user1), ErrUserExists)
	user2 := &entity.User{
		UUID:  uuid.NewString(),
		Login: "user2",
		Email: "user2@example.com",
	}
	require.NoError(t, userRepo.Create(ctx, user2))

	user3 := &entity.User{
		UUID:  uuid.NewString(),
		Login: "user3",
		Email: "user2@example.com",
	}
	require.NoError(t, userRepo.Create(ctx, user3))

	userDB, err := userRepo.GetByLogin(ctx, "www")
	require.ErrorContains(t, err, ErrUserNotExists)
	require.Nil(t, userDB)

	userDB, err = userRepo.GetByLogin(ctx, user1.Login)
	require.NoError(t, err)
	require.Equal(t, user1.UUID, userDB.UUID)
	require.Equal(t, user1.Email, userDB.Email)

	userDB, err = userRepo.GetByEmail(ctx, user2.Email)
	require.NoError(t, err)
	require.Equal(t, user2.UUID, userDB.UUID)
	require.Equal(t, user2.Login, userDB.Login)

	user2.Email = "user2new@example.com"
	user2.Hash = "hash"
	code := "code"
	user2.RegistrationCode = &code
	tm := time.Now().UTC().Round(time.Millisecond)
	user2.ConfirmedAt = &tm
	require.NoError(t, userRepo.Update(ctx, user2))

	userDB, err = userRepo.GetByUUID(ctx, user2.UUID)
	require.NoError(t, err)
	require.NotNil(t, userDB)
	require.Equal(t, user2.Email, userDB.Email)
	require.Equal(t, user2.Login, userDB.Login)
	require.Equal(t, user2.Hash, userDB.Hash)
	require.Equal(t, user2.RegistrationCode, userDB.RegistrationCode)
	require.Equal(t, user2.ConfirmedAt, userDB.ConfirmedAt)

	users, err := userRepo.GetExpired(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, []*entity.User{}, users)

	time.Sleep(time.Duration(2) * time.Second)

	users, err = userRepo.GetExpired(ctx, 1)
	require.NoError(t, err)
	require.Len(t, users, 1)

	user3.RegistrationCode = nil
	user3.ConfirmedAt = nil
	require.NoError(t, userRepo.Update(ctx, user3))

	userDB, err = userRepo.GetByUUID(ctx, user3.UUID)
	require.NoError(t, err)
	require.NotNil(t, userDB)
	require.Equal(t, user3.RegistrationCode, userDB.RegistrationCode)
	require.Equal(t, user3.ConfirmedAt, userDB.ConfirmedAt)

	users, err = userRepo.GetExpired(ctx, 1)
	require.NoError(t, err)
	require.Len(t, users, 1)

	user1.RegistrationCode = nil
	require.NoError(t, userRepo.Update(ctx, user1))

	userDB, err = userRepo.GetByUUID(ctx, user1.UUID)
	require.NoError(t, err)
	require.Equal(t, user1.RegistrationCode, userDB.RegistrationCode)

	require.NoError(t, userRepo.DeleteByUUID(ctx, user1.UUID))

	require.NoError(t, userRepo.DeleteByUUID(ctx, user3.UUID))

	_, err = userRepo.GetByUUID(ctx, user1.UUID)
	require.ErrorContains(t, err, ErrUserNotExists)
}
