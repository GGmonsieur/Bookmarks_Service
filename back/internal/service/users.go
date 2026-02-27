package service

import (
	"bookmark_service/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// api.POST("auth/register", svc.CreatUser)
func (s *Service) RegisterUser(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}
	err = s.UsersRepo.RegisterUserIndb(c.Request().Context(), &user)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(err.Error()))
	}

	return c.String(http.StatusOK, "Ok")
}

// api.POST("auth/login", svc.LoginUser)
func (s *Service) LoginUser(c echo.Context) error {
	var req *models.User
	if err := c.Bind(&req); err != nil {
		s.logger.Error(err)
		return c.JSON(http.StatusBadRequest, (InvalidParams))
	}

	// 2. Передаем в Repo только нужные данные
	// Изменяем сигнатуру Repo.Login, чтобы она принимала email и password
	token, err := s.UsersRepo.Login(c.Request().Context(), req)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	// 3. Возвращаем токен в JSON, чтобы фронтенд его сохранил
	return c.JSON(http.StatusOK, map[string]string{
		"access_token": token,
	})
}
