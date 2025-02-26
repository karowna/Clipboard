package database

import (
	"fmt"
	"os"
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