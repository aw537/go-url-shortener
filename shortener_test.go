package main

import (
	"strings"
	"testing"
)

func TestGenerateShortURL(t *testing.T) {
	// Define a table of test cases
	testCases := []struct {
		name     string
		longURL  string
		expected int // the expected length of the short URL
	}{
		{"Standard URL", "https://example.com/path", 8},
		{"Empty String", "", 8},
		{"URL with Query", "https://example.com/path?query=parameter", 8},
		{"URL with Anchor", "https://example.com/path#anchor", 8},
		{"Very Long URL", strings.Repeat("a", 1000), 8},
		{"URL with Scheme Missing", "example.com/path", 8},
	}

	uniqueShortURLs := make(map[string]bool)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shortURL := generateShortURL(tc.longURL)

			if len(shortURL) != tc.expected {
				t.Errorf("generateShortURL(%q) got length %v, want length %v", tc.longURL, len(shortURL), tc.expected)
			}

			// Check for hash collisions
			if _, exists := uniqueShortURLs[shortURL]; exists {
				t.Errorf("generateShortURL(%q) generated a duplicate short URL: %v", tc.longURL, shortURL)
			} else {
				uniqueShortURLs[shortURL] = true
			}
		})
	}
}
