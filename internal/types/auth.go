package types

import "time"

// Session represents the session model
// swagger:model
type Session struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    string    `json:"user_id"`
	IsBlocked bool      `json:"is_blocked" gorm:"default:false"`
	ExpiresAt time.Time `json:"expires_at" gorm:"type:datetime(3)"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`

	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// AuthToken holds authentication token details with refresh token
// swagger:model
type AuthToken struct {
	AccessToken  string `json:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// AuthUser represents data stored in JWT token for authenticated user
type AuthUser struct {
	ID    string
	Name  string
	Email string
	Role  string
	// add more if needed
}
