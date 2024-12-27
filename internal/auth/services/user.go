package services

import (
	"database/sql"
	"errors"
	"ssr-metaverse/internal/database"
	"ssr-metaverse/internal/auth/entities"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(username, password string) (*entities.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO account (username, password) VALUES ($1, $2) RETURNING id, created_at`
	var user entities.User
	err = database.DB.QueryRow(query, username, string(hashedPassword)).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	user.Username = username
	return &user, nil
}

func GetUserByID(id int) (*entities.User, error) {
	var user entities.User
	query := `SELECT id, username, created_at FROM account WHERE id = $1`
	err := database.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("usuário não encontrado")
	}
	return &user, err
}

func UpdateUser(id int, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `UPDATE account SET username = $1, password = $2 WHERE id = $3`
	_, err = database.DB.Exec(query, username, hashedPassword, id)
	return err
}

func DeleteUser(id int) error {
	query := `DELETE FROM account WHERE id = $1`
	_, err := database.DB.Exec(query, id)
	return err
}

func Authenticate(username, password string) (*entities.User, error) {
	var user entities.User
	query := `SELECT id, username, password FROM account WHERE username = $1`
	err := database.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return nil, errors.New("usuário não encontrado")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("senha incorreta")
	}

	return &user, nil
}