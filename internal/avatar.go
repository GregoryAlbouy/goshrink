package internal

// Avatar represents an user's avatar in the application.
type Avatar struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Url    string `json:"url"`
}

// AvatarService represents a service for managing users avatars.
type AvatarService interface {
	FindByUserID(id int) (*Avatar, error)

	CreateOne(url string) error

	UpdateOne(url string) error
}
