package models

type ClipItem struct {
	ID      uint   `gorm:"primaryKey"`
	Content string `json:"content,omitempty"`
	Image   string `json:"image,omitempty"`
}
type TemplateData struct {
	ClipItems []ClipItem
}
