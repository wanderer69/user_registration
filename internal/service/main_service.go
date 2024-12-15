package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/wanderer69/user_registration/internal/config"
	"github.com/wanderer69/user_registration/internal/gateway"
	"github.com/wanderer69/user_registration/internal/tools/mail"
	authUsecase "github.com/wanderer69/user_registration/internal/workflow/auth"
	apiPublic "github.com/wanderer69/user_registration/pkg/api/public"
)

type MainService struct {
	userRepository userRepository
}

func NewMainService(
	userRepository userRepository,
) *MainService {
	return &MainService{
		userRepository: userRepository,
	}
}

func (ms *MainService) Execute(ctx context.Context, cnf config.Config) {
	logger := zap.L()

	mailService := mail.NewMailer(cnf.ConfigMailer)

	authUseCase := authUsecase.NewAuthOperations(ms.userRepository, mailService, cnf.ConfigAuth)

	e := echo.New()
	e.HideBanner = true

	apiPublic.RegisterHandlers(
		e,
		gateway.NewServer(
			authUseCase,
		),
	)

	e.Any("/health-check/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, nil)
	})

	go func() {
		url := fmt.Sprintf(":%d", cnf.AppPort)
		logger.Info("Start server", zap.String("url", url))
		e.Start(url)
	}()
}
