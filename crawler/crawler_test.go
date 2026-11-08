package crawler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestCrawlDownloadsPages verifies that multiple webpages are downloaded.
func TestCrawlDownloadsPages(t *testing.T) {
	// Create a local HTTP test server.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		switch r.URL.Path {
		case "/robotics":
			fmt.Fprint(w, "<html><body>Robotics article</body></html>")
		case "/artificial-intelligence":
			fmt.Fprint(w, "<html><body>Artificial intelligence article</body></html>")
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	urls := []string{
		server.URL + "/robotics",
		server.URL + "/artificial-intelligence",
	}

	pages, err := Crawl(urls)

	if err != nil {
		t.Fatalf("expected successful crawl, got error: %v", err)
	}

	if len(pages) != 2 {
		t.Fatalf("expected 2 downloaded pages, got %d", len(pages))
	}

	// Verify that the crawler preserves the original URL order.
	for index, expectedURL := range urls {
		if pages[index].URL != expectedURL {
			t.Errorf(
				"page %d: expected URL %s, got %s",
				index,
				expectedURL,
				pages[index].URL,
			)
		}
	}

	// Verify that the downloaded HTML contains the expected content.
	if !strings.Contains(string(pages[0].HTML), "Robotics article") {
		t.Error("first page does not contain expected HTML content")
	}

	if !strings.Contains(string(pages[1].HTML), "Artificial intelligence article") {
		t.Error("second page does not contain expected HTML content")
	}
}

// TestCrawlEmptyURLs verifies that empty input returns no pages or errors.
func TestCrawlEmptyURLs(t *testing.T) {
	pages, err := Crawl(nil)

	if err != nil {
		t.Fatalf("expected no error for empty input, got: %v", err)
	}

	if len(pages) != 0 {
		t.Errorf("expected 0 pages, got %d", len(pages))
	}
}

// TestCrawlInvalidURL verifies that invalid URLs are reported as errors.
func TestCrawlInvalidURL(t *testing.T) {
	pages, err := Crawl([]string{"://invalid-url"})

	if err == nil {
		t.Fatal("expected an error for an invalid URL")
	}

	if len(pages) != 0 {
		t.Errorf("expected 0 pages, got %d", len(pages))
	}
}
