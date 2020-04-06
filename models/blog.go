package models

import (
	"encoding/json"
	"time"
)

// Blog represents data from Blogger API v3
type Blog struct {
	Kind        string    `json:"kind"`
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Published   time.Time `json:"published"`
	Updated     time.Time `json:"updated"`
	URL         string    `json:"url"`
	SelfLink    string    `json:"selfLink"`
	Posts       Posts     `json:"posts"`
	Pages       Pages     `json:"pages"`
	Locale      Locale    `json:"locale"`
}

// Locale represents locale of Blog of Blogger API v3
type Locale struct {
	Language string `json:"language"`
	Country  string `json:"country"`
	Variant  string `json:"variant"`
}

// Pages represents summary of pages from Blog of Blogger API v3
type Pages struct {
	TotalItems int    `json:"totalItems"`
	SelfLink   string `json:"selfLink"`
}

// Posts represents summary of posts from Blog of Blogger API v3
type Posts struct {
	TotalItems int    `json:"totalItems"`
	SelfLink   string `json:"selfLink"`
}

func (blog *Blog) String() string {
	b, err := json.Marshal(blog)
	if err != nil {
		return ""
	}
	return string(b)
}
