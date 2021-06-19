package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns a hash of the given password.
// Automatically and randomly generate and store a salt is upon hashing.
func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// HashPassword compares a hash to its possible plain text equivalent.
func ComparePasswords(hashedPwd string, plainPwd string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		return err
	}
	return nil
}
