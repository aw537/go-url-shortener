package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
)

var (
	// A map to store the URL mappings, with a mutex for concurrent access
	urlMap   = make(map[string]string)
	mapMutex = &sync.Mutex{}
)

func generateShortURL(longURL string) string {
	// SHA256 to hash the URL and return the first 8 characters
	hasher := sha256.New()
	hasher.Write([]byte(longURL))
	shortURL := hex.EncodeToString(hasher.Sum(nil))[:8]
	return shortURL
}

func shortenURLHandler(writer http.ResponseWriter, request *http.Request) {
	// Parse the long URL from the request
	longURL := request.URL.Query().Get("url")
	if longURL == "" {
		http.Error(writer, "Missing URL parameter", http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL(longURL)

	// Lock the map for safe concurrent access
	mapMutex.Lock()
	urlMap[shortURL] = longURL
	mapMutex.Unlock()

	// Respond with the short URL
	fmt.Fprintf(writer, "Short URL is: %s", shortURL)
}

func main() {
	http.HandleFunc("/shorten", shortenURLHandler)
	http.ListenAndServe(":8080", nil)
}
