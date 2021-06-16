package internal

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// User represents a user in the application.
type User struct {
	ID        int    `db:"id"         json:"id"`
	Username  string `db:"username"   json:"username"`
	Email     string `db:"email"      json:"email"`
	Password  string `db:"password"   json:"password"`
	AvatarURL string `db:"avatar_url" json:"avatar_url"`
}

// UserService represents a service for managing users.
type UserService interface {
	// FindByID retrieves a user by its ID.
	FindByID(id int) (User, error)

	// SetAvatarURL associates to given avatar URL to the given user ID.
	SetAvatarURL(userID int, url string) error

	// InsertOne inserts a user in the database.
	InsertOne(u User) error

	// Migrate inserts the given users in the database
	Migrate(users []User) error
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.EmailFormat),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 50)),
		validation.Field(&u.Username, validation.Required, is.Alphanumeric, validation.Length(3, 20)),
	)
}
