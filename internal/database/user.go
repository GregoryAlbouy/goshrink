package database

import (
	"github.com/GregoryAlbouy/shrinker/internal"
	"github.com/jmoiron/sqlx"
)

// Ensure service implements interface.
var _ internal.UserService = (*userService)(nil)

// UserService represents a service for managing users, it implements the interface internal.UserService.
type userService struct {
	db *sqlx.DB
}

// NewUserService returns a `userService` that implements
// `internal.UserService` interface.
func NewUserService(db *DB) internal.UserService {
	return &userService{db: db.sqlx}
}

// FindByID retrieves a user by its ID.
func (s *userService) FindByID(userID int) (internal.User, error) {
	u := internal.User{}

	if err := s.db.Get(
		&u,
		"SELECT * FROM V_user_avatar WHERE id = ?",
		userID,
	); err != nil {
		return internal.User{}, err
	}

	return u, nil
}

// FindCreds retrieves a user credentials by its username.
func (s *userService) FindCreds(username string) (internal.User, error) {
	u := internal.User{}

	if err := s.db.Get(
		&u,
		"SELECT id, username, password FROM user WHERE username = ?",
		username,
	); err != nil {
		return internal.User{}, err
	}

	return u, nil
}

func (s *userService) SetAvatarURL(userID int, url string) error {
	if _, err := s.db.Exec(
		"REPLACE INTO avatar (user_id, avatar_url) VALUES (?, ?)",
		userID, url,
	); err != nil {
		return err
	}

	return nil
}

// InsertOne inserts a user in the database.
func (s *userService) InsertOne(u internal.User) error {
	// Insert user
	if _, err := s.db.Exec(
		"INSERT INTO user (username, email, password) VALUES (?, ?, ?)",
		u.Username, u.Email, u.Password,
	); err != nil {
		return err
	}

	// Set avatar URL if they have one
	if u.AvatarURL != "" {
		return s.SetAvatarURL(u.ID, u.AvatarURL)
	}

	return nil
}

// Migrate inserts the given users in the database.
func (s *userService) Migrate(users []*internal.User) error {
	for _, u := range users {
		if err := s.InsertOne(*u); err != nil {
			return err
		}
	}

	return nil
}
