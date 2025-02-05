package services

import (
	"database/sql"
	"ssr-metaverse/internal/core/entities"
	"ssr-metaverse/internal/database"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/mock"
)

func TestAuthenticate(t *testing.T) {
	mockDB := new(database.MockDB)
	service := UserService{DB: mockDB}

	username := "testuser"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Configurar o mock para retornar o usuário correto
	mockDB.On(
		"QueryRow",
		"SELECT id_user, username, password FROM users WHERE username = $1",
		mock.Anything,
	).Return(&entities.User{ID: 1, Username: username, Password: string(hashedPassword)}, nil)

	t.Run("Autenticação bem-sucedida", func(t *testing.T) {
		user, err := service.Authenticate(username, password)
		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}
		if user.Username != username {
			t.Errorf("Esperava username '%s', mas recebeu '%s'", username, user.Username)
		}
	})

	t.Run("Usuário não encontrado", func(t *testing.T) {
		mockDB.On(
			"QueryRow",
			"SELECT id_user, username, password FROM users WHERE username = $1",
			"nonexistent",
		).Return(nil, sql.ErrNoRows)

		_, err := service.Authenticate("nonexistent", password)
		if err == nil || err.Error() != "usuário não encontrado" {
			t.Errorf("Esperava erro de usuário não encontrado, mas recebeu: %v", err)
		}
	})

	mockDB.AssertExpectations(t)
}