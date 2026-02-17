package main

import (
	"context"
	"bookmark_sevice/internal/config"
	"bookmark_sevice/internal/handlers"
	"bookmark_sevice/pkg/logs"
	"bookmark_sevice/pkg/postgres"

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

	svc := service.NewService(db, logger)

	router := echo.New()

	api := router.Group("/api/v1")

	api.POST("/bookmarks", svc.CreatBookmark)
	api.GET("/bookmarks/:id", svc.GetBookmarkFromID)
	api.GET("/bookmarks", svc.GETbookmarksPL)
	api.PATCH("/bookmarks/:id", svc.PATCHbookmarkid)
	api.DELETE("/bookmarks/:id", svc.DELETEid)
	router.Logger.Fatal(router.Start("localhost:" + cfg.GetWebPort()))
}
