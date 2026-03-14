package DataFunctions

import (
	"bookmark_service/internal/models"
	"context"
	"fmt"
)

// сoздание тега
func (r *Repo) InsertTag(ctx context.Context, tag *models.Tag) error {
	err := r.db.QueryRow(ctx, `INSERT INTO tags (name, user_id) VALUES ($1,$2) RETURNING id, created_at`,
		tag.Name, tag.UserID).
		Scan(&tag.ID, &tag.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// выдача тегов
func (r *Repo) FetchTags(ctx context.Context, userID int, page, limit int, search string) ([]models.Tag, error) {
	// Значения по умолчанию
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Базовый запрос
	sql := `SELECT id, user_id, name, created_at FROM tags WHERE user_id = $1`
	args := []interface{}{userID}
	argCount := 2

	// Поиск по названию тега
	if search != "" {
		sql += fmt.Sprintf(" AND name ILIKE $%d", argCount)
		args = append(args, "%"+search+"%")
		argCount++
	}

	// Сортировка по умолчанию (по имени)
	sql += fmt.Sprintf(" ORDER BY name ASC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	// Выполняем запрос
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var t models.Tag
		if err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return tags, rows.Err()
}

// удаление тега
func (s *Repo) DelTag(ctx context.Context, id, userID int) error {
	query := `DELETE FROM tags WHERE id = $1 and user_id = $2`

	res, err := s.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	// Дополнительная проверка: была ли вообще такая запись?
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("tag with id %d not found", id)
	}

	return nil
}
