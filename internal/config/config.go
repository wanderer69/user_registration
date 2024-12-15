package config

import (
	auth "github.com/wanderer69/user_registration/internal/workflow/auth"

	"github.com/wanderer69/user_registration/internal/tools/dao"
	"github.com/wanderer69/user_registration/internal/tools/mail"
)

type Config struct {
	AppPort uint   `envconfig:"APP_PORT" default:"8888"`
	AppEnv  string `envconfig:"APP_ENV" default:"prod"`

	dao.ConfigDAO
	mail.ConfigMailer
	auth.ConfigAuth
}
