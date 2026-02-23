package service

import (
	"net/http"
	"bookmark_service/internal/models"
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
    
    userID, ok := c.Get("user_id").(int)
    if !ok {
        s.logger.Error("user_id not found in context")
        return c.JSON(s.NewError(Unauthorized)) 
    }
    
    bookmark.UserID = userID
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
    userID, ok := c.Get("user_id").(int)
    if !ok {
        s.logger.Error("user_id not found in context")
        return c.JSON(s.NewError(Unauthorized)) 
    }

	repo := s.BookmarksRepo
    
	report, err := repo.GETbkmID(c.Request().Context(), id, userID)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.JSON(http.StatusOK, Response{Object: report})
}



//api.GET("/bookmarks", svc.GETbookmarksPL)
func (s *Service) GetBookmarksSort(c echo.Context) error {
    userID := c.Get("user_id").(int)

    filter := models.BookmarkFilter{
        Search: c.QueryParam("search"),
        TagID:  getOptionalInt(c.QueryParam("tagId")),
        Page:   getOptionalInt(c.QueryParam("page")),
        Limit:  getOptionalInt(c.QueryParam("limit")),
        Sort:   c.QueryParam("sort"),
        Order:  c.QueryParam("order"),
    }

    bookmarks, err := s.BookmarksRepo.FetchBookmarks(c.Request().Context(), userID, filter)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, err.Error())
    }

    return c.JSON(http.StatusOK, bookmarks)
}

func getOptionalInt(s string) int {
    val, _ := strconv.Atoi(s)
    return val
}

// api.PATCH("/bookmarks:id", svc.PATCHid)
func (s *Service) PATCHbookmarkid(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err!= nil{
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "неправильный парсер"})
    }

    

    var req models.UpdateBookmarkReq
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid json"})
    }

    userID, ok := c.Get("user_id").(int)
    if !ok {
        s.logger.Error("user_id not found in context")
        return c.JSON(s.NewError(Unauthorized)) 
    }

    err = s.BookmarksRepo.PatchBookmark(c.Request().Context(), id, userID, req.Title, req.Description)
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

    userID, ok := c.Get("user_id").(int)
    if !ok {
        s.logger.Error("user_id not found in context")
        return c.JSON(s.NewError(Unauthorized)) 
    } 

    err = s.BookmarksRepo.DeleteBookmark(c.Request().Context(), id,userID)
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