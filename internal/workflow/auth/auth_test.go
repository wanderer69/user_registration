package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRegistration(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepository := NewMockuserRepository(ctrl)
	mailService := NewMockmailService(ctrl)
	cnf := ConfigAuth{}
	ao := NewAuthOperations(userRepository, mailService, cnf)
	ctx := context.Background()
	require.ErrorContains(t, ao.Registration(ctx, "werwer"), "exported: password has non latin characters")
	require.ErrorContains(t, ao.Registration(ctx, "werwer"), "exported: password has lower lenght")
}
