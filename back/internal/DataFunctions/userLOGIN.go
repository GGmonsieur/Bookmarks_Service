package DataFunctions

import (
	"context"
	"errors"
	"time"
    "fmt"
    "github.com/jackc/pgx/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"bookmark_service/internal/models"
	"bookmark_service/internal/config"
)



// регистрация user
func (r *Repo) RegisterUserIndb(ctx context.Context, user *models.User) error {
    // 1. Хэшируем пароль
    // Cost 10-12 — оптимальный баланс между скоростью и безопасностью
    hashedBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.HashedPassword = string(hashedBytes)

    // 2. Выполняем запрос к БД
    query := `
        INSERT INTO users (email, hashed_password)
        VALUES ($1, $2)
        RETURNING id
    `
    
    err = r.db.QueryRow(ctx, query, user.Email, user.HashedPassword).
        Scan(&user.ID)
        
    if err != nil {
    if errors.Is(err, pgx.ErrNoRows) {
        return fmt.Errorf("user with email %s already exists", user.Email)
    }
    }
    return nil
}


//логин user
func (r *Repo) Login(ctx context.Context, user *models.User)(string, error) {
    
    // 1. Ищем пользователя в базе
    query := `SELECT id, hashed_password FROM users WHERE email = $1`
    err := r.db.QueryRow(ctx, query, &user.Email).Scan(&user.ID, &user.HashedPassword)
    if err != nil {
        
        return "", errors.New("invalid email or password")
    }

    // 2. Сравниваем хэш из базы с паролем из запроса
    err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(user.Password))
    if err != nil {
        return "", errors.New("invalid email or password")
    }

    // 3. Пароль верный! Создаем JWT токен
    claims := &models.MyCustomClaims{
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