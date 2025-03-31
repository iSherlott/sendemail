package models

type Attachment struct {
	Filename string `json:"filename"`
	Content  string `json:"content,omitempty"`
	URL      string `json:"url,omitempty"`
}
