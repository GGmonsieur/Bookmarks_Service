package bookmarks

import (
	"bookmark_sevice/internal/DataFunctions"
	"bookmark_sevice/pkg/postgres"

	"github.com/labstack/echo/v4"
)

const (
	InvalidParams       = "invalid params"
	InternalServerError = "internal error"
)

type Service struct {
	db     *postgres.DB
	logger echo.Logger

	Datas *DataFunctions.Repo
}

func NewData(db *postgres.DB, logger echo.Logger) *Service {
	svc := &Service{
		db:     db,
		logger: logger,
	}
	svc.initRepositories(db)

	return svc
}

func (s *Service) initRepositories(db *postgres.DB) {
	s.Datas = DataFunctions.NewRepo(db)
}

// Пока можно не вдаваться в то что ниже

type Response struct {
	Object       any    `json:"object,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
}

func (r *Response) Error() string {
	return r.ErrorMessage
}

func (s *Service) NewError(err string) (int, *Response) {
	return 400, &Response{ErrorMessage: err}
}
