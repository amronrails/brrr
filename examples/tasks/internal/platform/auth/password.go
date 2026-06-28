package auth

import "golang.org/x/crypto/bcrypt"

// bcryptCost balances security and latency; 12 is a sensible default for an
// interactive login flow.
const bcryptCost = 12

// HashPassword returns the bcrypt hash of a plaintext password.
func HashPassword(plaintext string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// CheckPassword reports whether plaintext matches the stored bcrypt hash.
func CheckPassword(hash, plaintext string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext)) == nil
}
