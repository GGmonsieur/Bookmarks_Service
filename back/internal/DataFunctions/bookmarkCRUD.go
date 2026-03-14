package DataFunctions

import (
	"bookmark_service/internal/models"
	"bookmark_service/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

type Repo struct {
	db *postgres.DB
}

func NewRepo(db *postgres.DB) *Repo {
	return &Repo{db: db}
}

// получение  bookmark через id
func (r *Repo) GETbkmID(ctx context.Context, id, userID int) (*models.Bookmark, error) {
	var bookmark models.Bookmark
	err := r.db.QueryRow(ctx, `SELECT id, user_id, url, title, description FROM bookmarks WHERE id = $1 AND user_id = $2 `, id, userID).
		Scan(&bookmark.ID, &bookmark.UserID, &bookmark.Url, &bookmark.Title, &bookmark.Description)
	if err != nil {
		return nil, err
	}
	return &bookmark, nil
}

// создание bookmar
func (r *Repo) POSTbookmark(ctx context.Context, bookmark *models.Bookmark) error {
	query := `INSERT INTO bookmarks (user_id, url, title, description) 
              VALUES ($1, $2, $3, $4) 
              RETURNING id, created_at`
	
	err := r.db.QueryRow(ctx, query, 
		bookmark.UserID, bookmark.Url, bookmark.Title, bookmark.Description,
	).Scan(&bookmark.ID, &bookmark.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { 
				return models.ErrDuplicateURL
			}
		}
		return err
	}

	return nil
}

// получение страниц и кол-во id
func (r *Repo) FetchBookmarks(ctx context.Context, userID int, f models.BookmarkFilter) ([]models.Bookmark, error) {
	// Значения по умолчанию
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 {
		f.Limit = 10
	}
	offset := (f.Page - 1) * f.Limit

	sql := `SELECT id, user_id, url, title, description FROM bookmarks WHERE user_id = $1 `
	args := []interface{}{userID}
	argCount := 2

	if !f.IncludDelete {
        sql += " AND deleted_at IS NULL"
    }
	// Фильтр по поиску (Title или URL)
	if f.Search != "" {
		sql += fmt.Sprintf(" AND (title ILIKE $%d OR url ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+f.Search+"%")
		argCount++
	}

	// Фильтр по тегу (если есть таблица связей)
	if f.TagID > 0 {
		sql += fmt.Sprintf(" AND id IN (SELECT bookmark_id FROM bookmark_tags WHERE tag_id = $%d)", argCount)
		args = append(args, f.TagID)
		argCount++
	}

	// Сортировка (ВНИМАНИЕ: нельзя вставлять через аргументы $1, только напрямую,
	// но нужно валидировать список разрешенных полей!)
	allowedSort := map[string]bool{"created_at": true, "title": true}
	if !allowedSort[f.Sort] {
		f.Sort = "created_at"
	}

	if strings.ToLower(f.Order) != "asc" {
		f.Order = "desc"
	}

	sql += fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", f.Sort, f.Order, argCount, argCount+1)
	args = append(args, f.Limit, offset)

	// Выполнение запроса
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var bookmarks []models.Bookmark

	for rows.Next() {

		var b models.Bookmark

		if err := rows.Scan(&b.ID, &b.UserID, &b.Url, &b.Title, &b.Description); err != nil {
			return nil, err
		}

		bookmarks = append(bookmarks, b)

	}
	return bookmarks, rows.Err()
}

// обновление title и description
func (s *Repo) PatchBookmark(ctx context.Context, id int, userID int, title *string, desc *string) error {
	// Собираем запрос динамически
	query := "UPDATE bookmarks SET "
	args := []any{}
	argID := 1

	if title != nil {
		query += fmt.Sprintf("title = $%d, ", argID)
		args = append(args, *title)
		argID++
	}
	if desc != nil {
		query += fmt.Sprintf("description = $%d, ", argID)
		args = append(args, *desc)
		argID++
	}

	// Если ничего не прислали — просто выходим
	if len(args) == 0 {
		return nil
	}
	query += "updated_at = CURRENT_TIMESTAMP"

	query += fmt.Sprintf(" WHERE id = $%d and user_id = $%d", argID, argID+1)
	args = append(args, id)
	args = append(args, userID)

	_, err := s.db.Exec(ctx, query, args...)
	return err
}

// soft удаление по id
func (s *Repo) DeleteBookmark(ctx context.Context, id, userID int) error {
	query := `
        UPDATE bookmarks 
        SET deleted_at = NOW() 
        WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
    `
	res, err := s.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	// Дополнительная проверка: была ли вообще такая запись?
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("bookmark with id %d not found", id)
	}

	return nil
}


//восстановление
func (s *Repo) Restore(ctx context.Context, id, userID int) error {
	query := `
        UPDATE bookmarks 
        SET deleted_at = NULL 
        WHERE id = $1 AND user_id = $2 AND deleted_at IS NOT NULL
    `
	res, err := s.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	// Дополнительная проверка: была ли вообще такая запись?
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("bookmark with id %d not found", id)
	}

	return nil
}
