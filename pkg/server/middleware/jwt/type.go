package jwt

// TokenOutput represents the output of a token request
type TokenOutput struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

// TokenInput represents the input of a token request
type TokenInput struct {
	Type   string                 `json:"type"` // refresh_token or access_token
	Claims map[string]interface{} `json:"claims"`
}
