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
	Password  string `db:"password"   json:"-"`
	AvatarURL string `db:"avatar_url" json:"avatar_url"`
}

// UserInput represents a user as sent by a client request.
type UserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserService represents a service for managing users.
type UserService interface {
	// FindByID retrieves a user by its ID.
	FindByID(id int) (User, error)

	// SetAvatarURL assigns to a user's avatar field the given avatar URL.
	SetAvatarURL(userID int, url string) error

	// InsertOne inserts a user in the database.
	InsertOne(u User) error

	// Migrate inserts the given users in the database
	Migrate(users []User) error
}

func NewUser(in UserInput) *User {
	return &User{
		Username: in.Username,
		Email:    in.Email,
		Password: in.Password,
	}
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.EmailFormat),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 50)),
		validation.Field(&u.Username, validation.Required, is.Alphanumeric, validation.Length(3, 20)),
	)
}
