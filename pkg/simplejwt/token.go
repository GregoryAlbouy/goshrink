package simplejwt

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidClaims = errors.New("invalid token claims")
	ErrInvalidToken  = errors.New("invalid token")
)

type BaseClaims struct {
	jwt.StandardClaims
}

func NewClaims(id string, exp time.Time) BaseClaims {
	return BaseClaims{
		jwt.StandardClaims{
			Id:        id,
			ExpiresAt: exp.Unix(),
		},
	}
}

// newToken returns a new token given claims and an expiration time.
func newToken(claims BaseClaims) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

// signToken returns the signed string from the given token.
//
// It panics if a secret key was not provided beforehand via SetSecretKey.
func signToken(token *jwt.Token) (string, error) {
	panicMissingSecretKey()
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", ErrInvalidClaims
	}
	return tokenString, nil
}

// NewSignedToken returns a new token as a signed a string.
func NewSignedToken(claims BaseClaims) (string, error) {
	token := newToken(claims)
	tokenString, err := signToken(token)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifiedToken parses, validates and returns a token given a string.
// Claims are then accessible for validation and use.
//
// It panics if a secret key was not provided beforehand via SetSecretKey.
func VerifiedToken(tokenString string) (*jwt.Token, error) {
	panicMissingSecretKey()

	// This method will return an error if the token is invalid
	// or if the signature does not match.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	if !token.Valid {
		return nil, ErrInvalidClaims
	}
	return token, nil
}

func ClaimsId(token jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidClaims
	}
	// "jti" is the json tag of StandardClaims.Id
	if _, ok := claims["jti"]; ok {
		if str, ok := claims["jti"].(string); ok {
			return str, nil
		}
	}
	return "", ErrInvalidClaims
}
