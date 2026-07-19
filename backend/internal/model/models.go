package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

type Task struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
}
