package middleware

import (
	"net/http"
	"os"
	
)

var AuthUsername string
var AuthPassword string

func InitializeVariables() {
	AuthUsername = os.Getenv("USERNAME")
	AuthPassword = os.Getenv("PASSWORD")
}
func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || !CheckUsernameAndPassword(username, password) {
			w.Header().Set("WWW-Authenticate", `Basic realm="DT-Webhooks"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
func CheckUsernameAndPassword(username, password string) bool {
	return username == AuthUsername && password == AuthPassword
}
