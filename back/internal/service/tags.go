package service

import (
	"bookmark_service/internal/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// protected.POST("/tags", svc.CreatTag)
func (s *Service) CreatTag(c echo.Context) error {
	var tag models.Tag
	err := c.Bind(&tag)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}

	userID, ok := c.Get("user_id").(int)
	if !ok {
		s.logger.Error("user_id not found in context")
		return c.JSON(s.NewError(Unauthorized))
	}

	tag.UserID = userID
	repo := s.TagsRepo
	err = repo.InsertTag(c.Request().Context(), &tag)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.String(http.StatusOK, "Ok")
}

// protected.GET("/tags", svc.GetTags)
func (s *Service) GetTags(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		s.logger.Error("user_id not found in context")
		return c.JSON(s.NewError(Unauthorized))
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}

	search := c.QueryParam("search")

	repo := s.TagsRepo
	tags, err := repo.FetchTags(c.Request().Context(), userID, page, limit, search)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.JSON(http.StatusOK, Response{Object: tags})
}

// protected.DELETE("tags/:id", svc.DeleteTag)
func (s *Service) DeleteTag(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id format"})
	}

	userID, ok := c.Get("user_id").(int)
	if !ok {
		s.logger.Error("user_id not found in context")
		return c.JSON(s.NewError(Unauthorized))
	}

	err = s.TagsRepo.DelTag(c.Request().Context(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		s.logger.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}

	// Возвращаем 204 No Content (стандарт для успешного удаления)
	return c.JSON(http.StatusOK, map[string]string{"status": "succes"})
}
