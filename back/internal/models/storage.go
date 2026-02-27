package models

import (
	"time"
)

type User struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	HashedPassword string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Bookmark struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Url         string    `json:"url"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type UpdateBookmarkReq struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type Tag struct {
	ID        int       `json:"id"`
	UserID    int       `json:"users_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type BookmarkFilter struct {
	Search string
	TagID  int
	Page   int
	Limit  int
	Sort   string // "created_at" или "title"
	Order  string // "asc" или "desc"
}
