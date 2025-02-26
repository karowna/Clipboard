package database

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"encoding/base64"
	"gorm.io/gorm"

	"github.com/karowna/Clipboard/src/models"
)


func DbConnectionString() string {
	var host string
	var username string
	var password string

	host = os.Getenv("DB_HOST")
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")

	return fmt.Sprintf("host=%s user=%s password=%s", host, username, password)
}

func CreateClip(w http.ResponseWriter, r *http.Request, db *gorm.DB) {

	var clipItem struct {
		Content string `json:"content,omitempty"`
		Image   string `json:"image,omitempty"`
	}

	err := json.NewDecoder(r.Body).Decode(&clipItem)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var imagePath string
	if clipItem.Image != "" {
		imageData, err := base64.StdEncoding.DecodeString(clipItem.Image)
		if err != nil {
			http.Error(w, "Invalid image data", http.StatusBadRequest)
			return
		}

		imageFileName := fmt.Sprintf("%d.png", time.Now().UnixNano())

		fullImagePath := filepath.Join("src", "static", "images", imageFileName) 

		os.MkdirAll(filepath.Dir(fullImagePath), os.ModePerm)

		err = os.WriteFile(fullImagePath, imageData, 0644)
		if err != nil {
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
		// different path needed for frontend
		imagePath = fmt.Sprintf("static/images/%s", imageFileName)
	}

	newClipItem := models.ClipItem{
		Content: clipItem.Content,
		Image:   imagePath,
	}

	result := db.Create(&newClipItem)
	if result.Error != nil {
		log.Printf("Error saving to database: %v", result.Error)
		http.Error(w, "Error saving to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Clipboard content created successfully"))
}



func GetClippedContent(db *gorm.DB) []models.ClipItem {
    log.Println("Getting clipped content from database")

    var clipItems []models.ClipItem
    db.Find(&clipItems)

    return clipItems
}

func DeleteAll(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	imageDir := "src/static/images/"

	err := filepath.Walk(imageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			os.Remove(path)
		}
		return nil
	})

	if err != nil {
		http.Error(w, "Failed to delete images", http.StatusInternalServerError)
		return
	}

	var clipItemsBefore []models.ClipItem
	db.Find(&clipItemsBefore) 
	log.Println("Clip items before deletion:", clipItemsBefore)

	db.Exec("DELETE FROM clip_items")

	var clipItemsAfter []models.ClipItem
	db.Find(&clipItemsAfter)
	log.Println("Clip items after deletion:", clipItemsAfter)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All images and database entries deleted"))
}

