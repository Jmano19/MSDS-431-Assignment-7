package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Jmano19/MSDS-431-Assignment-7/crawler"
	"github.com/Jmano19/MSDS-431-Assignment-7/scraper"
)

// createBenchmarkPages generates sample HTML pages for testing.
func createBenchmarkPages(count int) []crawler.Page {
	var pages []crawler.Page

	var content strings.Builder

	content.WriteString(`
		<html>
		<body>
			<h1 id="firstHeading">Robotics</h1>
			<div id="mw-content-text">
				<div class="mw-parser-output">
	`)

	// Generate paragraphs to simulate realistic article content.
	for i := 0; i < 20; i++ {
		content.WriteString(fmt.Sprintf(
			"<p>Robotics is an interdisciplinary field involving "+
				"engineering, computer science, artificial intelligence, "+
				"and mechanical systems. Paragraph %d describes "+
				"robotic applications and autonomous systems.</p>",
			i,
		))
	}

	content.WriteString(`
				</div>
			</div>
		</body>
		</html>
	`)

	html := []byte(content.String())

	for i := 0; i < count; i++ {
		pages = append(pages, crawler.Page{
			URL:  fmt.Sprintf("https://example.com/article/%d", i),
			HTML: html,
		})
	}

	return pages
}

// TestScrapePages verifies concurrent extraction and error handling.
func TestScrapePages(t *testing.T) {
	validHTML := []byte(`
		<html>
		<body>
			<h1 id="firstHeading">Test Article</h1>
			<div id="mw-content-text">
				<div class="mw-parser-output">
					<p>This is test article content.</p>
				</div>
			</div>
		</body>
		</html>
	`)

	pages := []crawler.Page{
		{
			URL:  "https://example.com/first",
			HTML: validHTML,
		},
		{
			URL:  "https://example.com/invalid",
			HTML: nil,
		},
		{
			URL:  "https://example.com/third",
			HTML: validHTML,
		},
	}

	articles, failures := scrapePages(pages, 2)

	if len(articles) != 2 {
		t.Fatalf("expected 2 successful articles, got %d", len(articles))
	}

	if failures != 1 {
		t.Errorf("expected 1 failure, got %d", failures)
	}

	// Verify that successful results retain their original order.
	for _, article := range articles {
		if article.Title != "Test Article" {
			t.Errorf("unexpected article title: %s", article.Title)
		}
	}

	if articles[0].URL != pages[0].URL {
		t.Errorf("first article is out of order")
	}

	if articles[1].URL != pages[2].URL {
		t.Errorf("second article is out of order")
	}
}

// BenchmarkScrapeSequential measures sequential article extraction.
func BenchmarkScrapeSequential(b *testing.B) {
	pages := createBenchmarkPages(10)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, page := range pages {
			_, err := scraper.Extract(page.URL, page.HTML)

			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

// BenchmarkScrapeConcurrent measures concurrent article extraction.
func BenchmarkScrapeConcurrent(b *testing.B) {
	pages := createBenchmarkPages(10)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		articles, failures := scrapePages(pages, workerCount)

		if failures != 0 {
			b.Fatalf("expected no failures, got %d", failures)
		}

		if len(articles) != len(pages) {
			b.Fatalf(
				"expected %d articles, got %d",
				len(pages),
				len(articles),
			)
		}
	}
}
