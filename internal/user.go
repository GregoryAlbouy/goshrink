package internal

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
