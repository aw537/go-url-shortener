package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGenerateShortURL(t *testing.T) {

	testCases := []struct {
		name     string
		longURL  string
		expected int // the expected length of the short URL
	}{
		{"Valid URL", "https://example.com/path", 8},
		{"Empty URL", "", 8},
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

func TestShortenURLHandler(t *testing.T) {

	testCases := []struct {
		name           string
		urlParam       string
		expectedStatus int
	}{
		{"Valid URL", "https://example.com", http.StatusOK},
		{"Empty URL", "", http.StatusBadRequest},
		{"Invalid URL", "http://%zzz", http.StatusBadRequest},
		{"Long URL Input", strings.Repeat("a", 1000), http.StatusOK},
		{"Special Characters", "https://example.com/!$&'()*+,;=:@/", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// Create a request to pass to the handler.
			req, err := http.NewRequest("GET", "/shorten?url="+tc.urlParam, nil)
			if err != nil {
				t.Fatal(err)
			}

			// ResponseRecorder records the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(shortenURLHandler)

			handler.ServeHTTP(rr, req)

			// Verify status code is what we expect.
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			// Verify response body is what we expect.
			if tc.expectedStatus == http.StatusOK {
				expected := "Short URL is: "
				if rr.Body.String()[:len(expected)] != expected {
					t.Errorf("handler returned unexpected body: got %v want %v",
						rr.Body.String(), expected)
				}
			}
		})
	}
}
