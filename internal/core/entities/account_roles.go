package entities

import "time"

type AccountRole struct {
	ID         int        `json:"id" db:"id_account_roles"`
	UserID     int        `json:"user_id" db:"id_user"`
	RoleID     int        `json:"role_id" db:"id_role"`
	GrantedAt  time.Time  `json:"granted_at" db:"granted_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
}
