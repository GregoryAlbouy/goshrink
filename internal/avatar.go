package internal

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// AvatarInput represents a user's avatar as sent by a client request.
// Note: it is a temporary implementation as in the end it will be the
// static server that will deal with storing and serving avatar files.
type AvatarInput struct {
	URL string `json:"avatar_url"`
}

func (a AvatarInput) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.URL, is.URL),
	)
}
