// internal/user/model/user.go
package model

type User struct {
	UserID           int64  `json:"user_id"`
	Username         string `json:"username,omitempty"`
	FirstName        string `json:"first_name,omitempty"`
	JoinedAt         string `json:"joined_at"`
	LastActiveAt     string `json:"last_active_at"`
	CalculationCount int    `json:"calculation_count"`
}
