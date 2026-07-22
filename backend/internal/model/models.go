package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

type Task struct {
	ID             int       `json:"id"`
	UserID         int       `json:"userId"`
	Title          string    `json:"title"`
	Done           bool      `json:"done"`
	AttachmentURL  *string   `json:"attachmentUrl"`
	AttachmentName *string   `json:"attachmentName"`
	AttachmentSize *int64    `json:"attachmentSize"`
	MimeType       *string   `json:"mimeType"`
	ThumbnailURL   *string   `json:"thumbnailUrl"`
	CreatedAt      time.Time `json:"createdAt"`
}

type TaskJob struct {
	TaskID    int       `json:"task_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
