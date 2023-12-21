package types

// Memo represents the memo model
// swagger:model
type Memo struct {
	Base
	UserID  string `json:"user_id"`
	Content string `json:"memo" gorm:"type:text"`
}
