package entities

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	CreatedAt string `json:"created_at"`
}