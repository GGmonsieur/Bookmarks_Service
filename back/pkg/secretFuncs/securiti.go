package secretFuncs

import (
	"bookmark_service/internal/config"
	"bookmark_service/internal/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type MyCustomClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func ChekPassword(user *models.User) (string, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(user.Password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// 3. Пароль верный! Создаем JWT токен
	claims := &MyCustomClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Токен на сутки
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	cfg, _ := config.NewConfig()
	jwtSecret := cfg.JWTsecret

	// Подписываем токен нашим секретным ключом
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
