package service

import (
	"net/http"
	"bookmark_sevice/internal/models"
	"strconv"
	"strings"
	"github.com/labstack/echo/v4"
)



//api.POST("/bookmarks", svc.CreatBookmark)
func (s *Service) CreatBookmark(c echo.Context) error {
	var bookmark models.Bookmark
	err := c.Bind(&bookmark)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}

	repo := s.BookmarksRepo
	err = repo.POSTbookmark(c.Request().Context(), &bookmark)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.String(http.StatusOK, "Ok")
}



//api.GET("/bookmarks:id", svc.GetBookmarkFromID)
func (s *Service) GetBookmarkFromID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}

	repo := s.BookmarksRepo

	report, err := repo.GETbkmID(c.Request().Context(), id)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.JSON(http.StatusOK, Response{Object: report})
}



//api.GET("/bookmarks", svc.GETbookmarksPL)
func (s *Service) GETbookmarksPL(c echo.Context) error {
    page, _ := strconv.Atoi(c.QueryParam("page"))
    limit, _ := strconv.Atoi(c.QueryParam("limit"))


    bookmarks, err := s.BookmarksRepo.FetchBookmarks(c.Request().Context(), page, limit)
    if err != nil {
        s.logger.Error(err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
    }

    return c.JSON(http.StatusOK, bookmarks)
}



// api.PATCH("/bookmarks:id", svc.PATCHid)
func (s *Service) PATCHbookmarkid(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err!= nil{
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "неправильный парсер"})
    }

    

    var req models.Bookmark
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid json"})
    }


    err = s.BookmarksRepo.PatchBookmark(c.Request().Context(), id, &req.Title, &req.Description)
    if err != nil {
        s.logger.Error(err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
    }

    return c.JSON(http.StatusOK,map[string]string{"OK": "поля bookmark и title обновлены"})
}



//api.DELETE("/bookmarks:id", svc.DELETEid)
func (s *Service) DELETEid(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id format"})
    }

    
    err = s.BookmarksRepo.DeleteBookmark(c.Request().Context(), id)
    if err != nil {
        if strings.Contains(err.Error(), "not found") {
            return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
        }
        s.logger.Error(err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
    }

    // Возвращаем 204 No Content (стандарт для успешного удаления)
    return c.JSON(http.StatusOK, map[string]string{"OK": "bookmark удален"})
}