package entities

type User struct {
	ID        int    `json:"id" db:"id_user"`
	Username  string `json:"username" db:"username"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"-" db:"password"`
	CreatedAt string `json:"created_at" db:"created_at"`
}
