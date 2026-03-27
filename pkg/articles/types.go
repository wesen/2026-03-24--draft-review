package articles

import "time"

type Section struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	BodyMarkdown string `json:"bodyMarkdown"`
}

type SectionInput struct {
	ID           string `json:"id,omitempty"`
	Title        string `json:"title"`
	BodyMarkdown string `json:"bodyMarkdown"`
}

type CreateArticleInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Intro  string `json:"intro"`
}

type UpdateArticleInput struct {
	Title    *string         `json:"title,omitempty"`
	Author   *string         `json:"author,omitempty"`
	Intro    *string         `json:"intro,omitempty"`
	Status   *string         `json:"status,omitempty"`
	Sections *[]SectionInput `json:"sections,omitempty"`
}

type CreateVersionInput struct {
	Label      string  `json:"label,omitempty"`
	Intro      *string `json:"intro,omitempty"`
	AuthorNote *string `json:"authorNote,omitempty"`
}

type Article struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Version   string    `json:"version"`
	Status    string    `json:"status"`
	Intro     string    `json:"intro"`
	ShareURL  string    `json:"shareUrl,omitempty"`
	Sections  []Section `json:"sections"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
