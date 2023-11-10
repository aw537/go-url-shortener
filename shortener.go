package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	// A map to store the URL mappings, with a mutex for concurrent access
	urlMap   = make(map[string]string)
	mapMutex = &sync.Mutex{}
	// keeps track of the amount of times the URL has been accessed
	urlAccessCount = make(map[string]int)
)

func generateShortURL(longURL string) string {
	// SHA256 to hash the URL and return the first 8 characters
	hasher := sha256.New()
	hasher.Write([]byte(longURL))
	shortURL := hex.EncodeToString(hasher.Sum(nil))[:8]
	return shortURL
}

func shortenURLHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Method is not supported.", http.StatusNotFound)
		return
	}

	// Parse the long URL from the request
	longURL := request.URL.Query().Get("url")
	if longURL == "" {
		http.Error(writer, "Missing URL", http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL(longURL)

	// Lock the map for safe concurrent access
	mapMutex.Lock()
	urlMap[shortURL] = longURL
	mapMutex.Unlock()

	// Respond with the short URL
	fmt.Fprintf(writer, "The shortened URL is: %s", shortURL)
}

func redirectHandler(writer http.ResponseWriter, request *http.Request) {
	// Extract the code from the request path
	code := request.URL.Path[1:] // Remove the leading "/"

	// Look up the code in the URL map
	mapMutex.Lock()
	longURL, ok := urlMap[code]

	if !ok {
		// If the code is not found, return a 404 not found error
		http.NotFound(writer, request)
		return
	} else {
		urlAccessCount[code]++
	}
	mapMutex.Unlock()

	// Redirect to the original URL
	http.Redirect(writer, request, longURL, http.StatusFound)
}

func statsHandler(writer http.ResponseWriter, request *http.Request) {
	// Extract the code from the request path
	code := request.URL.Path[len("/stats/"):]

	// Retrieve and return the access count for the code
	mapMutex.Lock()
	count, ok := urlAccessCount[code]
	mapMutex.Unlock()

	if !ok {
		http.Error(writer, "Short URL not found", http.StatusNotFound)
		return
	}

	// Respond with the count
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(fmt.Sprintf("Access count for %s: %d", code, count)))
}

func main() {
	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/shorten", shortenURLHandler)
	http.HandleFunc("/stats/", statsHandler)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
