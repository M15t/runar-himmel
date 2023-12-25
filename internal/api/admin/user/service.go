package user

import (
	"runar-himmel/internal/repo"
	"runar-himmel/pkg/rbac"
)

// New creates new user application service
func New(repo *repo.Service, rbacSvc rbac.Intf, cr Crypter) *User {
	return &User{repo: repo, rbac: rbacSvc, cr: cr}
}

// User represents user application service
type User struct {
	repo *repo.Service
	rbac rbac.Intf
	cr   Crypter
}

// Crypter represents security interface
type Crypter interface {
	CompareHashAndPassword(string, string) bool
	HashPassword(string) string
}
