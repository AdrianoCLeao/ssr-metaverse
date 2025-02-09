package services

import (
	"database/sql"
	"ssr-metaverse/internal/core/entities"
	"ssr-metaverse/internal/database"
	"ssr-metaverse/internal/core/error"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DB database.DBInterface
}

func (s *UserService) CreateUser(username, email, password string) (*entities.User, *error.APIError) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, &error.APIError{
			Code:    500,
			Message: "Error generating password hash.",
		}
	}

	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id_user, created_at`
	var user entities.User
	err = s.DB.QueryRow(query, username, email, string(hashedPassword)).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code {
			case "23505":
				return nil, &error.APIError{
					Code:    409,
					Message: "Username or Email already taken",
				}
			default:
				return nil, &error.APIError{
					Code:    500,
					Message: "Unexpected error occurred creating User.",
				}
			}
		}

		return nil, &error.APIError{
			Code:    500,
			Message: "Unexpected error occurred creating User",
		}
	}

	return &user, nil
}

func (s *UserService) GetUserByID(id int) (*entities.User, *error.APIError) {
	var user entities.User
	query := `SELECT id_user, username, email, created_at FROM users WHERE id_user = $1`
	err := s.DB.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, &error.APIError{
			Code:    404,
			Message: "User not found",
		}
	} else if err != nil {
		return nil, &error.APIError{
			Code:    500,
			Message: "Unexpected error occurred fetching user",
		}
	}

	return &user, nil
}

func (s *UserService) UpdateUser(id int, username, password string) *error.APIError {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return &error.APIError{
			Code:    500,
			Message: "Error generating password hash.",
		}
	}

	query := `UPDATE users SET username = $1, password = $2 WHERE id_user = $3`
	result, err := s.DB.Exec(query, username, hashedPassword, id)
	if err != nil {
		return &error.APIError{
			Code:    500,
			Message: "Unexpected error occurred updating user",
		}
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return &error.APIError{
			Code:    404,
			Message: "User not found",
		}
	}

	return nil
}

func (s *UserService) DeleteUser(id int) *error.APIError {
	query := `DELETE FROM users WHERE id_user = $1`
	result, err := s.DB.Exec(query, id)
	if err != nil {
		return &error.APIError{
			Code:    500,
			Message: "Unexpected error occurred deleting user",
		}
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return &error.APIError{
			Code:    404,
			Message: "User not found",
		}
	}

	return nil
}

func (s *UserService) Authenticate(username, password string) (*entities.User, *error.APIError) {
	var user entities.User
	query := `SELECT id_user, username, password FROM users WHERE username = $1`
	err := s.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return nil, &error.APIError{
			Code:    404,
			Message: "User not found",
		}
	} else if err != nil {
		return nil, &error.APIError{
			Code:    500,
			Message: "Unexpected error occurred during authentication",
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, &error.APIError{
			Code:    401,
			Message: "Incorrect password",
		}
	}

	return &user, nil
}

func (s *UserService) GetUserRoles(userID int) ([]string, *error.APIError) {
	query := `
		SELECT r.role_name
		FROM account_roles ar
		JOIN roles r ON ar.id_role = r.id_role
		WHERE ar.id_user = $1 AND ar.revoked_at IS NULL
	`

	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, &error.APIError{
			Code:    500,
			Message: "Unexpected error occurred fetching user roles",
		}
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var roleName string
		if err := rows.Scan(&roleName); err != nil {
			return nil, &error.APIError{
				Code:    500,
				Message: "Unexpected error occurred scanning roles",
			}
		}
		roles = append(roles, roleName)
	}

	if len(roles) == 0 {
		return nil, &error.APIError{
			Code:    404,
			Message: "No roles found for user",
		}
	}

	return roles, nil
}
