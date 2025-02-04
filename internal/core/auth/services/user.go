package services

import (
	"database/sql"
	"errors"
	"ssr-metaverse/internal/core/entities"
	"ssr-metaverse/internal/database"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DB database.DBInterface
}

func (s *UserService) CreateUser(username, password string) (*entities.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id_user, created_at`
	var user entities.User
	err = s.DB.QueryRow(query, username, string(hashedPassword)).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	user.Username = username

	roleQuery := `
		INSERT INTO account_roles (id_user, id_role) 
		SELECT $1, id_role FROM roles WHERE role_name = 'user'
	`
	_, err = s.DB.Exec(roleQuery, user.ID)
	if err != nil {
		return nil, errors.New("erro ao atribuir o cargo 'user' ao usuário")
	}

	return &user, nil
}

func (s *UserService) GetUserByID(id int) (*entities.User, error) {
	var user entities.User
	query := `SELECT id_user, username, created_at FROM users WHERE id_user = $1`
	err := s.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("usuário não encontrado")
	}
	return &user, err
}

func (s *UserService) UpdateUser(id int, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `UPDATE users SET username = $1, password = $2 WHERE id_user = $3`
	_, err = s.DB.Exec(query, username, hashedPassword, id)
	return err
}

func (s *UserService) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id_user = $1`
	_, err := s.DB.Exec(query, id)
	return err
}

func (s *UserService) Authenticate(username, password string) (*entities.User, error) {
	var user entities.User
	query := `SELECT id_user, username, password FROM users WHERE username = $1`
	err := s.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return nil, errors.New("usuário não encontrado")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("senha incorreta")
	}

	return &user, nil
}

func (s *UserService) GetUserRoles(userID int) ([]string, error) {
	query := `
		SELECT r.role_name
		FROM account_roles ar
		JOIN roles r ON ar.id_role = r.id_role
		WHERE ar.id_user = $1 AND ar.revoked_at IS NULL
	`

	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var roleName string
		if err := rows.Scan(&roleName); err != nil {
			return nil, err
		}
		roles = append(roles, roleName)
	}

	if len(roles) == 0 {
		return nil, errors.New("nenhum cargo encontrado para o usuário")
	}

	return roles, nil
}
