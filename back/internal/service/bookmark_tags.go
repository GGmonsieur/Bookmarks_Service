package service

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// protected.POST("/bookmarks/:id/tags", svc.PostBKM_TAGS)
func (s *Service) PostBKM_TAGS(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	bookmarkID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid bookmark id format"})
	}

	var req struct {
		TagIDs []int `json:"tagIds"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid json body"})
	}

	// Проверка на пустой массив (опционально, но полезно)
	if len(req.TagIDs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "tagIds cannot be empty"})
	}

	ctx := c.Request().Context()
	err = s.BookmarksRepo.AddTagsToBookmark(ctx, bookmarkID, userID, req.TagIDs)

	if err != nil {
		s.logger.Error("failed to link tags: ", err)

		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "denied") || strings.Contains(err.Error(), "invalid or belong") {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Access denied: check if bookmark and all tags belong to you",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal database error"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"status": "tags successfully linked to bookmark"})
}

// protected.DELETE("/bookmarks/:id/tags/:tagId", svc.DeleteBKM_TAG)
func (s *Service) DeleteBKM_TAG(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	bookmarkID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid bookmark id"})
	}

	tagID, err := strconv.Atoi(c.Param("tagId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tag id"})
	}

	err = s.BookmarksRepo.RemoveTagFromBookmark(c.Request().Context(), bookmarkID, tagID, userID)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "succes"})
}
