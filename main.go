package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// map to store shortened URLs to the original ones
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

// handler to create a short URL
func shortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		originalURL := r.PostFormValue("url") // retrieve URL from form

		// check if the URL starts with http:// or https://
		if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
			originalURL = "http://" + originalURL
		}

		shortToken := generateShortToken()
		urlMap[shortToken] = originalURL
		fmt.Fprintf(w, "http://localhost:8080/r/%s", shortToken)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // initialize random generator

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Handle URL shortening
	http.HandleFunc("/shorten", shortenURL)

	// Handle redirection
	http.HandleFunc("/r/", redirectFromShort)

	log.Println("The server is now running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func redirectFromShort(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	
	if strings.HasPrefix(path, "/r/") {
		shortToken := path[len("/r/"):] 
		if originalURL, exists := urlMap[shortToken]; exists {
			http.Redirect(w, r, originalURL, http.StatusSeeOther)
			return
		} else {
			http.Error(w, "Short URL not found", http.StatusNotFound)
			return
		}
	}

}
