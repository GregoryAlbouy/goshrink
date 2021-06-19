package simplejwt

import "log"

var secretKey []byte

// SetSecretKey sets the secret key used for hashing JWT.
// It must be used before any usage of this package.
func SetSecretKey(key []byte) {
	secretKey = key
}

func panicMissingSecretKey() {
	if len(secretKey) == 0 {
		log.Panic("missing secret key: call jsonwebtoken.SetSecretKey before using jsonwebtoken module")
	}
}
