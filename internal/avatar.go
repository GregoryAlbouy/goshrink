package internal

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Avatar represents an user's avatar in the application.
type Avatar struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	URL    string `json:"url"`
}

// AvatarService represents a service for managing users avatars.
type AvatarService interface {
	FindByUserID(id int) (*Avatar, error)

	CreateOne(url string) error

	UpdateOne(url string) error
}

func (a Avatar) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.URL, is.URL),
	)
}
