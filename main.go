package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// Map to store shortened URLs to the original ones
var urlMap = make(map[string]string)

// Function to generate a shortened URL token
func generateShortToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6) // we'll use 6 characters for the shortened URL
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// Handler to create a short URL
func shortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		originalURL := r.PostFormValue("url") // Retrieve URL from form

		// Check if the URL starts with http:// or https://
		if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
			// If not, prepend "http://"
			originalURL = "http://" + originalURL
		}

		shortToken := generateShortToken()
		urlMap[shortToken] = originalURL
		fmt.Fprintf(w, "http://localhost:8080/%s", shortToken)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// Handler to redirect from the short URL to the original one
func redirectFromShort(w http.ResponseWriter, r *http.Request) {
	shortToken := r.URL.Path[len("/"):]
	if originalURL, exists := urlMap[shortToken]; exists {
		http.Redirect(w, r, originalURL, http.StatusSeeOther)
	} else {
		http.Error(w, "Short URL not found", http.StatusNotFound)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // initialize global pseudo random generator

	http.HandleFunc("/shorten", shortenURL) // Register the handler
	http.HandleFunc("/", redirectFromShort) // Assumes all other requests are for short URLs

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
