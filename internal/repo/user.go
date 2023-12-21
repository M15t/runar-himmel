package repo

import (
	"context"
	"runar-himmel/internal/types"

	repoutil "runar-himmel/pkg/util/repo"

	"gorm.io/gorm"
)

// User represents the client for user table
type User struct {
	*repoutil.Repo[types.User]
}

// NewUser returns a new user database instance
func NewUser(gdb *gorm.DB) *User {
	return &User{repoutil.NewRepo[types.User](gdb)}
}

// FindByEmail finds a user by the given email
func (r *User) FindByEmail(ctx context.Context, email string) (rec *types.User, err error) {
	rec = &types.User{}
	err = r.GDB.Take(rec, `email = ?`, email).Error

	return
}

// UpdateRefreshToken updates the refresh token of the given user
func (r *User) UpdateRefreshToken(ctx context.Context, userID, refreshToken string) error {
	return r.GDB.Model(&types.User{}).Where(`id = ?`, userID).Update(`refresh_token`, refreshToken).Error
}
