package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/wanderer69/user_registration/internal/entity"
)

type dao interface {
	DB() *gorm.DB
}

type Repository struct {
	db dao
}

type model struct {
	Uuid             string
	Login            string
	Email            string
	RegistrationCode *string
	Hash             string
	AccessCode       string
	CreatedAt        time.Time
	ConfirmedAt      *time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

type models []model

func convert(u *entity.User) (model, error) {
	m := model{
		Uuid:             u.UUID,
		Login:            u.Login,
		Email:            u.Email,
		RegistrationCode: u.RegistrationCode,
		Hash:             u.Hash,
		AccessCode:       u.AccessCode,
		ConfirmedAt:      u.ConfirmedAt,
	}
	return m, nil
}

func (m model) convert() (*entity.User, error) {
	v := &entity.User{
		UUID:             m.Uuid,
		Login:            m.Login,
		Email:            m.Email,
		RegistrationCode: m.RegistrationCode,
		Hash:             m.Hash,
		AccessCode:       m.AccessCode,
		CreatedAt:        m.CreatedAt,
		ConfirmedAt:      m.ConfirmedAt,
		UpdatedAt:        m.UpdatedAt,
		DeletedAt:        m.DeletedAt,
	}
	return v, nil
}

func NewRepository(db dao) *Repository {
	return &Repository{
		db: db,
	}
}

const (
	UsersTableName      = "users"
	ErrUserExists       = "user exists"
	ErrUserNotExists    = "user not exists"
	ErrResultQueryEmpty = "result query empty"
	uniqueViolation     = "23505"
)

func (r *Repository) Create(ctx context.Context, user *entity.User) error {
	m, err := convert(user)
	if err != nil {
		return err
	}
	m.CreatedAt = time.Now().UTC().Round(time.Millisecond)
	m.UpdatedAt = m.CreatedAt
	err = r.db.DB().WithContext(ctx).
		Table(UsersTableName).
		Create(m).Error
	if err != nil {
		if strings.Contains(err.Error(), uniqueViolation) {
			return errors.New(ErrUserExists)
		}
		return err
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, user *entity.User) error {
	m, err := convert(user)
	if err != nil {
		return err
	}
	m.UpdatedAt = time.Now().UTC().Round(time.Millisecond)
	err = r.db.DB().WithContext(ctx).
		Table(UsersTableName).
		Select("*").
		Where("uuid = ?", m.Uuid).
		Omit(
			"uuid",
			"created_at",
		).
		Updates(m).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) ConfirmationUpdate(ctx context.Context, user *entity.User) error {
	m, err := convert(user)
	if err != nil {
		return err
	}
	m.UpdatedAt = time.Now().UTC().Round(time.Millisecond)
	err = r.db.DB().WithContext(ctx).
		Table(UsersTableName).
		Select("*").
		Where("uuid = ?", m.Uuid).
		Omit(
			"uuid",
			"created_at",
		).
		Updates(m).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) selectOne(ctx context.Context, db *gorm.DB) (*entity.User, error) {
	var (
		m   model
		err error
	)

	tx := db.WithContext(ctx).Table(UsersTableName + " AS u").Where("u.deleted_at IS NULL")

	err = tx.Take(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(ErrUserNotExists)
		}
		return nil, err
	}

	return m.convert()
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return r.selectOne(ctx, r.db.DB().Where("u.email = ?", email))
}

func (r *Repository) GetByLogin(ctx context.Context, login string) (*entity.User, error) {
	return r.selectOne(ctx, r.db.DB().Where("u.login = ?", login))
}

func (r *Repository) GetByUUID(ctx context.Context, uuid string) (*entity.User, error) {
	return r.selectOne(ctx, r.db.DB().Where("u.uuid = ?", uuid))
}

func (r *Repository) GetByRegistrationCode(ctx context.Context, code string) (*entity.User, error) {
	return r.selectOne(ctx, r.db.DB().Where("u.registration_code = ?", code))
}

func (r *Repository) selectMany(ctx context.Context, db *gorm.DB) ([]*entity.User, error) {
	var (
		m   models
		err error
	)

	tx := db.WithContext(ctx).Table(UsersTableName + " AS u").Where("u.deleted_at IS NULL")

	err = tx.Find(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(ErrUserNotExists)
		}
		return nil, err
	}
	result := []*entity.User{}
	for i := range m {
		c, err := m[i].convert()
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}

	return result, nil
}

func (r *Repository) GetExpired(ctx context.Context, timeout int) ([]*entity.User, error) {
	return r.selectMany(ctx, r.db.DB().
		Where("EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - created_at)) > ?", timeout).
		Where("u.registration_code IS NOT NULL"))
}

func (r *Repository) DeleteByUUID(ctx context.Context, userUUID string) error {
	dt := time.Now().UTC().Round(time.Millisecond)
	m := model{
		Uuid:      userUUID,
		DeletedAt: &dt,
	}
	err := r.db.DB().WithContext(ctx).
		Table(UsersTableName).
		Select("*").
		Where("uuid = ?", m.Uuid).
		Where("deleted_at IS NULL").
		Omit(
			"uuid",
			"login",
			"email",
			"registration_code",
			"hash",
			"access_code",
			"confirmed_at",
		).
		Updates(m).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteByDeletedAt(ctx context.Context) error {
	m := model{}
	return r.db.DB().WithContext(ctx).
		Table(UsersTableName).
		Where("deleted_at IS NOT NULL").
		Delete(m).Error
}
