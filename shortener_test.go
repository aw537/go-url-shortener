package main

import (
	"fmt"
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
				expected := "The shortened URL is: "
				if rr.Body.String()[:len(expected)] != expected {
					t.Errorf("handler returned unexpected body: got %v want %v",
						rr.Body.String(), expected)
				}
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	// Initialize the URL map with some data for testing.
	urlMap = map[string]string{
		"abc123": "http://example.com",
		"xyz789": "http://example.org",
	}

	// Define test cases
	testCases := []struct {
		name         string
		shortCode    string
		expectedURL  string
		expectedCode int
	}{
		{"Valid Redirect", "abc123", "http://example.com", http.StatusFound},
		{"Non-existent Code", "doesNotExist", "", http.StatusNotFound},
		{"Special Characters", "a!b@c#d$", "", http.StatusNotFound},
		{"Case Sensitivity", "ABC123", "", http.StatusNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new HTTP request to the short URL
			req := httptest.NewRequest("GET", "/"+tc.shortCode, nil)
			// Create a ResponseRecorder to record the response
			w := httptest.NewRecorder()

			// Call the handler function
			redirectHandler(w, req)

			// Check the status code
			if w.Code != tc.expectedCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedCode, w.Code)
			}

			// If a redirect is expected, check the Location header
			if tc.expectedCode == http.StatusFound {
				location, err := w.Result().Location()
				if err != nil {
					t.Fatal(err)
				}
				if location.String() != tc.expectedURL {
					t.Errorf("Expected redirect to %q, got %q", tc.expectedURL, location.String())
				}
			}
		})
	}
}

func TestStatsHandler(t *testing.T) {

	urlAccessCount["new123"] = 0
	urlAccessCount["abc123"] = 3
	urlAccessCount["xyz789"] = 10

	// Define test cases
	testCases := []struct {
		name         string
		shortCode    string
		expectedBody string
		expectedCode int
	}{
		{"Existing Short URL", "abc123", "Access count for abc123: 3", http.StatusOK},
		{"Non-Existent Short URL", "nope123", "Short URL not found\n", http.StatusNotFound},
		{"Zero Accesses", "new123", "Access count for new123: 0", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/stats/%s", tc.shortCode), nil)
			w := httptest.NewRecorder()

			statsHandler(w, req)

			// Verify the status code
			if w.Code != tc.expectedCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedCode, w.Code)
			}

			// Verify the body of the response
			if w.Body.String() != tc.expectedBody {
				t.Errorf("Expected body %q, got %q", tc.expectedBody, w.Body.String())
			}
		})
	}
}
