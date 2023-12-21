package types

import (
	"runar-himmel/pkg/util/ulidutil"
	"time"

	"gorm.io/gorm"
)

// Base contains common fields for all models
type Base struct {
	// ID of the record
	ID string `json:"id" gorm:"primaryKey"`
	// The time that record is created
	CreatedAt time.Time `json:"created_at"`
	// The latest time that record is updated
	UpdatedAt time.Time `json:"updated_at"`
	// Enable soft-delete
	// DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// BeforeCreate hook executed by gorm
func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = ulidutil.NewString()
	}
	return
}
