package repo

import (
	"context"
	"runar-himmel/internal/types"

	repoutil "runar-himmel/pkg/util/repo"

	"gorm.io/gorm"
)

// Session represents the client for session table
type Session struct {
	*repoutil.Repo[types.Session]
}

// NewSession returns a new session database instance
func NewSession(gdb *gorm.DB) *Session {
	return &Session{repoutil.NewRepo[types.Session](gdb)}
}

// FindByID finds a session by the given ID and preload User
func (r *Session) FindByID(ctx context.Context, id, userID string) (rec *types.Session, err error) {
	rec = &types.Session{}
	err = r.GDB.Preload(`User`).Take(rec).Where(`id = ? AND user_id = ? AND is_blocked = false`, id, userID).Error

	return
}

// DeleteExpired deletes expired sessions
func (r *Session) DeleteExpired(ctx context.Context, userID string) error {
	return r.GDB.Delete(&types.Session{}, `expires_at < NOW() AND user_id = ?`, userID).Error
}
