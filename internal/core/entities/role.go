package entities

type Role struct {
	ID          int    `json:"id" db:"id_role"`
	RoleName    string `json:"role_name" db:"role_name"`
	Description string `json:"description,omitempty" db:"description"`
}
