package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/wanderer69/user_registration/internal/entity"
)

type ConfigAuth struct {
	InternalAdminUserName     string `envconfig:"INTERNAL_ADMIN_USER_NAME" default:"admin"`
	InternalAdminUserPassword string `envconfig:"INTERNAL_ADMIN_USER_PASSWORD" default:"adminpassword"`
	InternalAdminIPMask       string `envconfig:"INTERNAL_ADMIN_IP_MASK" default:"127.0.0.1"`
	RegistrationSubject       string `envconfig:"MAIL_REGISTRATION_SUBJ" default:"Registration"`
	MessageFormat             string `envconfig:"MAIL_REGISTRATION_FORMAT" default:"Dear client!\r\nRegistration code %v\r\n"`
	FromName                  string `envconfig:"INTERNAL_ADMIN_IP_MASK" default:"Admin"`
}

type AuthOperations struct {
	userRepository userRepository
	mailService    mailService
	cnf            ConfigAuth
}

func NewAuthOperations(
	userRepository userRepository,
	mailService mailService,
	cnf ConfigAuth,
) *AuthOperations {
	return &AuthOperations{
		userRepository: userRepository,
		mailService:    mailService,
		cnf:            cnf,
	}
}

const (
	restoreSubject          string = "Restoring access to the services personal account"
	UserPermissionRoleAdmin string = "admin"
)

func CreatePasswordHash(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (uo *AuthOperations) Registration(ctx context.Context, email string) error {
	logger := zap.L()
	logger.Info("Registration", zap.String("email", email))
	// проверяем, что нет такого почтового ящика
	_, err := uo.userRepository.GetByEmail(ctx, email)
	if err == nil {
		return fmt.Errorf("exported: found user with email %v", email)
	}

	registrationSubject := uo.cnf.RegistrationSubject
	messageFormat := uo.cnf.MessageFormat
	fromName := uo.cnf.FromName

	// отправляем по почте код
	code := uuid.NewString()[0:6]
	msg := fmt.Sprintf(messageFormat, code)
	err = uo.mailService.Send(email, registrationSubject, msg, fromName)
	if err != nil {
		return fmt.Errorf("exported: failed sent email %v: %v", email, err)
	}
	// сохраняем пользователя и код
	user := &entity.User{
		Email: email,
	}

	user.UUID = uuid.NewString()
	user.RegistrationCode = &code

	return uo.userRepository.Create(ctx, user)
}

func (uo *AuthOperations) ConfirmationOTP(ctx context.Context, code string, email string) error {
	logger := zap.L()
	logger.Info("ConfirmationOTP", zap.String("code", code), zap.String("email", email))
	// проверяем, что такой код еще есть
	user, err := uo.userRepository.GetByRegistrationCode(ctx, code)
	if err != nil {
		return fmt.Errorf("exported: not found user with registration code %v", code)
	}
	if user.Email != email {
		return fmt.Errorf("user with code %v not have email %v", code, email)
	}
	user.RegistrationCode = nil

	return uo.userRepository.ConfirmationUpdate(ctx, user)
}

func (uo *AuthOperations) Confirmation(ctx context.Context, email string, password string, login string) error {
	logger := zap.L()
	logger.Info("Confirmation", zap.String("email", email), zap.String("login", login), zap.String("password", password))
	passwordLen := int64(8)

	// проверить что символы только латинские
	passwordCnt := 0
	for i, w := 0, 0; i < len(password); i += w {
		runeValue, width := utf8.DecodeRuneInString(password[i:])
		if unicode.IsLetter(runeValue) {
			if !unicode.Is(unicode.Latin, runeValue) {
				return fmt.Errorf("password has non latin characters")
			}
		}
		w = width
		passwordCnt += 1
	}
	// проверить минимальную длину пароля
	if passwordCnt < int(passwordLen) {
		return fmt.Errorf("password has lower lenght")
	}

	// проверяем, что нет такого логина
	_, err := uo.userRepository.GetByLogin(ctx, login)
	if err == nil {
		return fmt.Errorf("found user with login %v", login)
	}

	// проверяем, что такой код еще есть
	user, err := uo.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("not found user with email %v", email)
	}

	if user.RegistrationCode != nil {
		return fmt.Errorf("user with email %v not confirmed", email)
	}

	if len(login) == 0 {
		return fmt.Errorf("login must be not empty")
	}
	_, err = uo.userRepository.GetByLogin(ctx, login)
	if err == nil {
		return fmt.Errorf("found user with login %v", login)
	}

	user.Hash = CreatePasswordHash(password)
	t := time.Now().UTC().Round(time.Millisecond)
	user.ConfirmedAt = &t
	user.Login = login

	return uo.userRepository.ConfirmationUpdate(ctx, user)
}

func (uo *AuthOperations) Login(ctx context.Context, login string, password string) (string, error) {
	user, err := uo.userRepository.GetByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("not found user with login %v", login)
	}

	hash := CreatePasswordHash(password)

	if user.RegistrationCode != nil && len(*user.RegistrationCode) > 0 {
		return "", fmt.Errorf("found not confirmed user with login %v", login)
	}
	if hash != user.Hash {
		return "", fmt.Errorf("bad password or login %v", login)
	}

	code := uuid.NewString()

	err = uo.userRepository.Update(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed update user with uuid %v", user.UUID)
	}
	return code, nil
}
