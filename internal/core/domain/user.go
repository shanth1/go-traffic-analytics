package domain

import "time"

type UserID string
type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleClient UserRole = "client"
)

type User struct {
	ID                 UserID     `json:"id"`
	Email              string     `json:"email"`
	PasswordHash       string     `json:"password_hash"`
	Role               UserRole   `json:"role"`
	PlanID             string     `json:"plan_id"`
	PlanExpires        *time.Time `json:"plan_expires"` // nil = forever (for free)
	ClicksCurrentMonth int        `json:"clicks_current_month"`
	IsActive           bool       `json:"is_active"`
	CreatedAt          time.Time  `json:"created_at"`
	DeletedAt          *time.Time `json:"-"`
}

func (u *User) ToPublic() UserPublic {
	return UserPublic{
		ID:                 u.ID,
		Email:              u.Email,
		Role:               string(u.Role),
		PlanID:             u.PlanID,
		PlanExpires:        u.PlanExpires,
		ClicksCurrentMonth: u.ClicksCurrentMonth,
		IsActive:           u.IsActive,
		CreatedAt:          u.CreatedAt,
	}
}

type UserPublic struct {
	ID                 UserID     `json:"id"`
	Email              string     `json:"email"`
	Role               string     `json:"role"`
	PlanID             string     `json:"plan_id"`
	PlanExpires        *time.Time `json:"plan_expires"`
	ClicksCurrentMonth int        `json:"clicks_current_month"`
	IsActive           bool       `json:"is_active"`
	CreatedAt          time.Time  `json:"created_at"`
}
