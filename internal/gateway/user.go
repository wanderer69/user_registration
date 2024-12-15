package gateway

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	api "github.com/wanderer69/user_registration/pkg/api/public"
)

const (
	ResultOK    = "Ok"
	ResultError = "Error"
)

// (GET /api/v1/public/user/confirmation)
func (s *Server) Confirmation(c echo.Context, params api.ConfirmationParams) error {
	logger := zap.L()

	ctx := c.Request().Context()
	logger.Info("HTTP.Confirmation handler started")
	err := s.authOperations.Confirmation(ctx, params.Email, params.Password, params.Login)
	if err != nil {
		errMsg := err.Error()
		return c.JSON(http.StatusOK, &api.OperationResponses{
			Result:       ResultError,
			ErrorMessage: &errMsg,
		})
	}
	return c.JSON(http.StatusOK, &api.OperationResponses{
		Result: ResultOK,
	})
}

// (GET /api/v1/public/user/confirmation_otp)
func (s *Server) ConfirmationOTP(c echo.Context, params api.ConfirmationOTPParams) error {
	logger := zap.L()

	ctx := c.Request().Context()
	logger.Info("HTTP.ConfirmationOTP handler started")
	err := s.authOperations.ConfirmationOTP(ctx, params.Otp, params.Email)
	if err != nil {
		errMsg := err.Error()
		return c.JSON(http.StatusOK, &api.OperationResponses{
			Result:       ResultError,
			ErrorMessage: &errMsg,
		})
	}
	return c.JSON(http.StatusOK, &api.OperationResponses{
		Result: ResultOK,
	})
}

// (GET /api/v1/public/user/register)
func (s *Server) Register(c echo.Context, params api.RegisterParams) error {
	logger := zap.L()

	ctx := c.Request().Context()
	logger.Info("HTTP.Register handler started")
	err := s.authOperations.Registration(ctx, params.Email)
	if err != nil {
		errMsg := err.Error()
		return c.JSON(http.StatusOK, &api.OperationResponses{
			Result:       ResultError,
			ErrorMessage: &errMsg,
		})
	}
	return c.JSON(http.StatusOK, &api.OperationResponses{
		Result: ResultOK,
	})
}

// (GET /api/v1/public/user/login)
func (s *Server) Login(c echo.Context, params api.LoginParams) error {
	logger := zap.L()

	ctx := c.Request().Context()
	logger.Info("HTTP.Login handler started")
	code, err := s.authOperations.Login(ctx, params.Login, params.Password)
	if err != nil {
		errMsg := err.Error()
		return c.JSON(http.StatusOK, &api.LoginResponses{
			Code:         "",
			ErrorMessage: &errMsg,
		})
	}
	return c.JSON(http.StatusOK, &api.LoginResponses{
		Code: code,
	})
}
