package DataFunctions

import (
	"context"
	"fmt"
	"bookmark_sevice/internal/models"
	"bookmark_sevice/pkg/postgres"
)

type Repo struct {
	db *postgres.DB
}

func NewRepo(db *postgres.DB) *Repo {
	return &Repo{db: db}
}



// получение  bookmark через id
func (r *Repo) GETbkmID(ctx context.Context, id int) (*models.Bookmark, error) {
	var bookmark models.Bookmark
	err := r.db.QueryRow(ctx, `SELECT id, user_id, url, title, description FROM bookmarks WHERE id = $1`, id).
		Scan(&bookmark.ID,&bookmark.UserID, &bookmark.Url, &bookmark.Title, &bookmark.Description)
	if err != nil {
		return nil, err
	}
	return &bookmark, nil
}



// создание bookmark
func (r *Repo) POSTbookmark(ctx context.Context, bookmark *models.Bookmark) error {
	err := r.db.QueryRow(ctx, `INSERT INTO bookmarks (user_id, url, title, description) VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		bookmark.UserID, bookmark.Url, bookmark.Title, bookmark.Description).
		Scan(&bookmark.ID, &bookmark.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}



//получение страниц и кол-во id
func (s *Repo) FetchBookmarks(ctx context.Context, page, limit int) ([]models.Bookmark, error) {
    if page < 1 { page = 1 }
    if limit < 1 { limit = 10 }
    offset := (page - 1) * limit

    query := `SELECT id, url, title, description FROM bookmarks LIMIT $1 OFFSET $2`
    rows, err := s.db.Query(ctx, query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var bookmarks []models.Bookmark
    for rows.Next() {
        var b models.Bookmark
        if err := rows.Scan(&b.ID, &b.Url, &b.Title, &b.Description); err != nil {
            return nil, err
        }
        bookmarks = append(bookmarks, b)
    }

    return bookmarks, rows.Err()
}



// обновление title и description
func (s *Repo) PatchBookmark(ctx context.Context, id int, title *string, desc *string) error {
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
    
    query += fmt.Sprintf(" WHERE id = $%d", argID)
    args = append(args, id)

    _, err := s.db.Exec(ctx, query, args...)
    return err
}



// удаление по id
func (s *Repo) DeleteBookmark(ctx context.Context, id int) error {
    query := `DELETE FROM bookmarks WHERE id = $1`


    res, err := s.db.Exec(ctx, query, id)
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