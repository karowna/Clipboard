package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "runtime"

    "github.com/gorilla/mux"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "github.com/karowna/Clipboard/src/database"
    "github.com/karowna/Clipboard/src/models"
    "github.com/karowna/Clipboard/src/middleware"
)

var indexTemplate *template.Template

func init() {
    _, filename, _, ok := runtime.Caller(0)
    if !ok {
        log.Fatal("Failed to determine current file's path")
    }

    templatePath := filepath.Join(filepath.Dir(filename), "models" , "index.html")

    var err error
    indexTemplate, err = template.ParseFiles(templatePath)
    if err != nil {
        log.Fatal("Error parsing template: ", err)
    }
}

func main() {
	
    middleware.InitializeVariables()

    db, err := gorm.Open(postgres.Open(database.DbConnectionString()), &gorm.Config{})
    if err != nil {
        log.Panic(err)
    }
    db.AutoMigrate(&models.ClipItem{})

    router := mux.NewRouter()

    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("src/static"))))

	router.HandleFunc("/", middleware.BasicAuth(func(w http.ResponseWriter, r *http.Request) {
		clippedContent := database.GetClippedContent(db)
	
		data := models.TemplateData{
			ClipItems: clippedContent,
		}
	
		err := indexTemplate.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))

    router.HandleFunc("/paste", middleware.BasicAuth(func(w http.ResponseWriter, r *http.Request) {
        database.CreateClip(w, r, db)
    })).Methods("POST")
	
	router.HandleFunc("/delete-all", func(w http.ResponseWriter, r *http.Request) {
		database.DeleteAll(w, r, db)
	})

    port := os.Getenv("PORT")
    if port == "" {
        port = "8088"
    }
    hostString := fmt.Sprintf("0.0.0.0:%s", port)
    log.Printf("Starting server on %s", hostString)

    if err := http.ListenAndServe(hostString, router); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
