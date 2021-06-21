package simplejwt

var secretKey []byte

// SetSecretKey sets the secret key used for hashing JWT.
// It must be used before any usage of this package.
func SetSecretKey(key []byte) {
	secretKey = key
}
