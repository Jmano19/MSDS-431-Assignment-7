package scraper

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/Jmano19/MSDS-431-Assignment-7/models"
)

// Extract converts downloaded HTML into a structured article.
func Extract(url string, html []byte) (models.Article, error) {
	if len(html) == 0 {
		return models.Article{}, fmt.Errorf("empty HTML content for %s", url)
	}

	// Parse the downloaded HTML document.
	document, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return models.Article{}, fmt.Errorf("parse HTML for %s: %w", url, err)
	}

	// Extract the article title.
	title := strings.TrimSpace(document.Find("#firstHeading").First().Text())
	if title == "" {
		title = strings.TrimSpace(document.Find("title").First().Text())
	}

	if title == "" {
		return models.Article{}, fmt.Errorf("article title not found for %s", url)
	}

	// Locate the main Wikipedia article content.
	content := document.Find("#mw-content-text .mw-parser-output").First()

	if content.Length() == 0 {
		return models.Article{}, fmt.Errorf("article content not found for %s", url)
	}

	// Remove HTML elements that are not part of the article's main text.
	content.Find(
		"script, style, noscript, nav, table, " +
			".navbox, .metadata, .reference, .reflist, " +
			".infobox, .sidebar, .mw-editsection, " +
			".hatnote, .toc, .catlinks, .shortdescription",
	).Remove()

	// Extract readable paragraphs, headings, and list items.
	var textParts []string

	content.Find("p, h2, h3, h4, li").Each(func(_ int, selection *goquery.Selection) {
		text := strings.Join(strings.Fields(selection.Text()), " ")

		if text != "" {
			textParts = append(textParts, text)
		}
	})

	articleText := strings.Join(textParts, "\n\n")

	if strings.TrimSpace(articleText) == "" {
		return models.Article{}, fmt.Errorf("no article text extracted from %s", url)
	}

	// Return the structured article.
	return models.Article{
		Title:     title,
		URL:       url,
		Text:      articleText,
		ScrapedAt: time.Now().UTC(),
	}, nil
}
