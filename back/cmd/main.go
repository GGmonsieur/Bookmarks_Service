package main

import (
	"bookmark_service/internal/config"
	"bookmark_service/internal/middlewear"
	"bookmark_service/internal/service"
	"bookmark_service/pkg/logs"
	"bookmark_service/pkg/postgres"
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	ctx := context.Background()
	defer ctx.Done()

	logger := logs.NewLogger(false)

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal(err)
	}

	db, err := postgres.Connect(ctx, cfg.Postgres)
	if err != nil {
		logger.Fatal(err)
	}
	log.Info("Postgres successfully connected")

	svc := service.NewData(db, logger)

	router := echo.New()

	v1 := router.Group("/api/v1")

	v1.POST("/auth/register", svc.RegisterUser)
	v1.POST("/auth/login", svc.LoginUser)

	protected := v1.Group("")
	protected.Use(middlewear.JWTMiddleware(cfg.JWTsecret))

	protected.POST("/bookmarks", svc.CreatBookmark)
	protected.GET("/bookmarks/:id", svc.GetBookmarkFromID)
	protected.GET("/bookmarks", svc.GetBookmarksSort)
	protected.PATCH("/bookmarks/:id", svc.PATCHbookmarkid)
	protected.DELETE("/bookmarks/:id", svc.DELETEid)

	protected.POST("/tags", svc.CreatTag)
	protected.GET("/tags", svc.GetTags)
	protected.DELETE("/tags/:id", svc.DeleteTag)

	protected.POST("/bookmarks/:id/tags", svc.PostBKM_TAGS)
	protected.DELETE("/bookmarks/:id/tags/:tagId", svc.DeleteBKM_TAG)
	router.Logger.Fatal(router.Start("localhost:" + cfg.GetWebPort()))
}
