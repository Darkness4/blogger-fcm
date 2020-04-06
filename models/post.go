package models

import "time"

// Post model from Blogger API v3
type Post struct {
	Kind      string    `json:"kind"`
	ID        string    `json:"id"`
	Published time.Time `json:"published"`
	Updated   time.Time `json:"updated"`
	URL       string    `json:"url"`
	SelfLink  string    `json:"selfLink"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}
