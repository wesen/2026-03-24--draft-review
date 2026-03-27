package articleassets

import "time"

type Asset struct {
	ID               string    `json:"id"`
	ArticleID        string    `json:"articleId"`
	URL              string    `json:"url"`
	Markdown         string    `json:"markdown"`
	OriginalFilename string    `json:"originalFilename"`
	ContentType      string    `json:"contentType"`
	ByteSize         int64     `json:"byteSize"`
	CreatedAt        time.Time `json:"createdAt"`
}

type UploadInput struct {
	Filename string
	Content  []byte
}

type createAssetRecordInput struct {
	ID               string
	StorageKey       string
	OriginalFilename string
	ContentType      string
	ByteSize         int64
}

type storedAsset struct {
	ID               string
	ArticleID        string
	StorageKey       string
	OriginalFilename string
	ContentType      string
	ByteSize         int64
	CreatedAt        time.Time
}
