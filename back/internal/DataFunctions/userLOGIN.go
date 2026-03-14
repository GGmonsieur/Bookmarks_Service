package DataFunctions

import (
	"context"
	"errors"
	"fmt"

    "bookmark_service/internal/models"
	"bookmark_service/pkg/secretFuncs"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
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
    tokenString,err := secretFuncs.ChekPassword(user)
    if err != nil {
        return "", errors.New("invalid email or password")
    }

    return tokenString, nil
}