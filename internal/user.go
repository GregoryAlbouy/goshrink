package internal

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// User represents a user in the application.
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	AvatarURL string `json:"avatar_url"`
}

// UserService represents a service for managing users.
type UserService interface {
	FindByID(id int) (*User, error)

	CreateOne(u *User) error

	UpdateOne(u *User) error
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.EmailFormat),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 50)),
		validation.Field(&u.Username, validation.Required, is.Alphanumeric, validation.Length(3, 20)),
	)
}
