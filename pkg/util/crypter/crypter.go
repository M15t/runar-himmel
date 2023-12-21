package crypter

import (
	"golang.org/x/crypto/bcrypt"
)

// New initalizes crypter service
func New() *Service {
	return &Service{}
}

// Service holds crypter methods
type Service struct{}

// HashPassword hashes the password using bcrypt
func (*Service) HashPassword(password string) string {
	return HashPassword(password)
}

// CompareHashAndPassword matches hash with password. Returns true if hash and password match.
func (*Service) CompareHashAndPassword(hash, password string) bool {
	return CompareHashAndPassword(hash, password)
}

///// Static functions /////

// HashPassword hashes the password using bcrypt
func HashPassword(password string) string {
	hashedPW, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPW)
}

// CompareHashAndPassword matches hash with password. Returns true if hash and password match.
func CompareHashAndPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
