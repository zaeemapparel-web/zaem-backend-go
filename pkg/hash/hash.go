package hash

import "golang.org/x/crypto/bcrypt"

const Cost = 12

// Password hashes a plain password
func Password(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), Cost)
	return string(bytes), err
}

// Compare checks if plain matches hash
func Compare(hashed, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}