package internal

import (
	"fmt"
	"net/url"
)

// Avatar represents an user's avatar in the application.
type Avatar struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	Url    string `json:"url"`
}

// AvatarService represents a service for managing users avatars.
type AvatarService interface {
	FindByUserID(id uint) (*Avatar, error)

	CreateOne(url string) error

	UpdateOne(url string) error
}

func NewAvatarURL(blob string) string {
	base := &url.URL{
		Scheme: "http",
		Host:   "company",
	}
	return fmt.Sprintf("%s/%s", base.String(), blob)
}
