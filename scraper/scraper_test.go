package scraper

import (
	"strings"
	"testing"
)

// TestExtractValidHTML verifies successful article extraction.
func TestExtractValidHTML(t *testing.T) {
	html := []byte(`
		<html>
			<head><title>Robotics - Wikipedia</title></head>
			<body>
				<h1 id="firstHeading">Robotics</h1>
				<div id="mw-content-text">
					<div class="mw-parser-output">
						<p>Robotics is an interdisciplinary field.</p>
						<h2>History</h2>
						<p>Robotics has developed over many years.</p>
						<ul>
							<li>Industrial robots</li>
							<li>Service robots</li>
						</ul>
						<nav>Navigation content</nav>
						<table><tr><td>Table content</td></tr></table>
						<div class="infobox">Infobox content</div>
					</div>
				</div>
			</body>
		</html>
	`)

	article, err := Extract("https://en.wikipedia.org/wiki/Robotics", html)
	if err != nil {
		t.Fatalf("expected successful extraction, got error: %v", err)
	}

	if article.Title != "Robotics" {
		t.Errorf("expected title Robotics, got %s", article.Title)
	}

	if article.URL != "https://en.wikipedia.org/wiki/Robotics" {
		t.Errorf("unexpected URL: %s", article.URL)
	}

	if !strings.Contains(article.Text, "interdisciplinary field") {
		t.Error("article text does not contain expected paragraph")
	}

	if !strings.Contains(article.Text, "History") {
		t.Error("article text does not contain expected heading")
	}

	if !strings.Contains(article.Text, "Industrial robots") {
		t.Error("article text does not contain expected list item")
	}

	if strings.Contains(article.Text, "Navigation content") {
		t.Error("article text should not contain navigation content")
	}

	if strings.Contains(article.Text, "Table content") {
		t.Error("article text should not contain table content")
	}

	if strings.Contains(article.Text, "Infobox content") {
		t.Error("article text should not contain infobox content")
	}

	if article.ScrapedAt.IsZero() {
		t.Error("expected a valid scrape timestamp")
	}
}

// TestExtractEmptyHTML verifies that empty HTML returns an error.
func TestExtractEmptyHTML(t *testing.T) {
	_, err := Extract("https://en.wikipedia.org/wiki/Robotics", nil)

	if err == nil {
		t.Error("expected an error for empty HTML")
	}
}

// TestExtractMissingTitle verifies that missing titles return an error.
func TestExtractMissingTitle(t *testing.T) {
	html := []byte(`
		<html>
			<body>
				<div id="mw-content-text">
					<div class="mw-parser-output">
						<p>Some article content.</p>
					</div>
				</div>
			</body>
		</html>
	`)

	_, err := Extract("https://example.com/article", html)

	if err == nil {
		t.Error("expected an error when article title is missing")
	}
}

// TestExtractMissingContent verifies that missing article content returns an error.
func TestExtractMissingContent(t *testing.T) {
	html := []byte(`
		<html>
			<body>
				<h1 id="firstHeading">Robotics</h1>
				<p>Content outside the article container.</p>
			</body>
		</html>
	`)

	_, err := Extract("https://example.com/article", html)

	if err == nil {
		t.Error("expected an error when article content is missing")
	}
}

// TestExtractNoReadableText verifies that empty article content returns an error.
func TestExtractNoReadableText(t *testing.T) {
	html := []byte(`
		<html>
			<body>
				<h1 id="firstHeading">Robotics</h1>
				<div id="mw-content-text">
					<div class="mw-parser-output">
						<nav>Navigation only</nav>
						<table><tr><td>Table only</td></tr></table>
					</div>
				</div>
			</body>
		</html>
	`)

	_, err := Extract("https://example.com/article", html)

	if err == nil {
		t.Error("expected an error when no readable article text exists")
	}
}
