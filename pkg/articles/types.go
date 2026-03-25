package articles

import "time"

type Section struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Paragraphs []string `json:"paragraphs"`
}

type Article struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Version   string    `json:"version"`
	Status    string    `json:"status"`
	Intro     string    `json:"intro"`
	Sections  []Section `json:"sections"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
