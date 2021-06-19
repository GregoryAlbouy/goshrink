package crypto_test

import (
	"math/rand"
	"testing"

	"github.com/GregoryAlbouy/shrinker/pkg/crypto"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" +
	"0124356789" +
	"~`!@#$%^&*()_-+={[}]|\\:;\"'<,>.?/"

func randCharSeq(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func TestComparePasswords(t *testing.T) {
	var passwords = [10][]byte{}
	for _, pwd := range passwords {
		pwd = []byte(randCharSeq(32))

		pwdString := string(pwd)
		hash, err := crypto.HashPassword(pwdString)
		if err != nil {
			t.Errorf("could not hash \"%s\"", pwdString)
		}

		if err = crypto.ComparePasswords(hash, pwdString); err != nil {
			t.Errorf("\"%s\" and \"%s\" should be equivalent", hash, pwdString)
		}
	}
}
