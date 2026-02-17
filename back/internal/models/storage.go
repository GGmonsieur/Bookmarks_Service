package models

import "time"

type User struct {
	ID             int64     `json:"id"`
	Email          string    `json:"email"`
	HashedPassword int       `json:"hashed_password"`
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

type Tag struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"users_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
