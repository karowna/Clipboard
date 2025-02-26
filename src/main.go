package main

import (
	"log"
	"runtime"
	"path/filepath"
	"html/template"
    "net/http"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"github.com/gorilla/mux"

	"github.com/karowna/Clipboard/src/database"
)

var indexTemplate *template.Template

func init() {
	// Build filepath for index.html page
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Failed to determine current file's path")
	}
	templatePath := filepath.Join(filepath.Dir(filename), "models", "index.html")
	var err error
	indexTemplate, err = template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal("Error parsing template: ", err)
	}
}

func main() {
	
	db, err := gorm.Open(postgres.Open(database.DbConnectionString()), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Panic(err)
	}
	db.AutoMigrate(&models.Webhook{})

	router := mux.NewRouter()

	router.PathPrefix("/index/static/").Handler(http.StripPrefix("/index/static/", http.FileServer(http.Dir("src/static"))))
}