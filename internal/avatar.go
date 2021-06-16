package internal

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// AvatarInput represents a user's avatar as sent by a client request.
type AvatarInput struct {
	URL string `json:"avatar_url"`
}

func (a AvatarInput) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.URL, is.URL),
	)
}
