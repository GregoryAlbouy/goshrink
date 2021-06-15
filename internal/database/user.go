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
//
// FIXME: this methods returns an error if the user does not have an avatar.
// This is due to the assignement of NULL to a string. See sql.NullString.
// Still looking for a fix at the moment.
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

// Migrate inserts the given users in the database
func (s *userService) Migrate(users []internal.User) error {
	for _, u := range users {
		if err := s.InsertOne(u); err != nil {
			return err
		}
	}

	return nil
}
