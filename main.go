package main

import (
	"crypto/sha256"
	"encoding/hex"
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
