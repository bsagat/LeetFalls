package models

import "time"

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	ReplyToID int       `json:"reply_to_id,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	ImageLink string    `json:"image_link"`
	Author    User      `json:"author"`
}
