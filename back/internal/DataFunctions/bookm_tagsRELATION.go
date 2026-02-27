package DataFunctions

import (
	"context"
	"fmt"
)

// создание связи
func (r *Repo) AddTagsToBookmark(ctx context.Context, bookmarkID int, userID int, tagIDs []int) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var bookmarkExists bool
	err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM bookmarks WHERE id = $1 AND user_id = $2)`,
		bookmarkID, userID).Scan(&bookmarkExists)
	if err != nil {
		return err
	}
	if !bookmarkExists {
		return fmt.Errorf("bookmark not found or access denied")
	}

	var validTagsCount int
	err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM tags WHERE id = ANY($1) AND user_id = $2`,
		tagIDs, userID).Scan(&validTagsCount)

	if err != nil {
		return err
	}
	if validTagsCount != len(tagIDs) {
		return fmt.Errorf("one or more tags are invalid or belong to another user")
	}

	// 3. Если обе проверки прошли — делаем вставку
	query := `INSERT INTO bookmark_tags (bookmark_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	for _, tagID := range tagIDs {
		if _, err := tx.Exec(ctx, query, bookmarkID, tagID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// удаление
func (r *Repo) RemoveTagFromBookmark(ctx context.Context, bookmarkID, tagID, userID int) error {
	query := `
        DELETE FROM bookmark_tags 
        WHERE bookmark_id = $1 
          AND tag_id = $2
          AND bookmark_id IN (SELECT id FROM bookmarks WHERE user_id = $3)
          AND tag_id IN (SELECT id FROM tags WHERE user_id = $4) `

	res, err := r.db.Exec(ctx, query, bookmarkID, tagID, userID, userID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении связи: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("связь не найдена или доступ запрещен")
	}

	return nil
}
