package entity

import (
	"time"
)

type User struct {
	UUID             string     `json:"uuid"`
	Login            string     `json:"login"`
	Email            string     `json:"email"`
	RegistrationCode *string    `json:"registration_code"`
	Hash             string     `json:"hash"`
	AccessCode       string     `json:"access_code"`
	ConfirmedAt      *time.Time `json:"confirmed_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
}
