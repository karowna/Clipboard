package models

type ClipItem struct {
	ID           uint
	Content string `json:"content,omitempty"`
    Image   []byte `json:"image,omitempty"`
}