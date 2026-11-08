package models

import "time"

// Article represents the text content extracted from a webpage.
type Article struct {
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Text      string    `json:"text"`
	ScrapedAt time.Time `json:"scraped_at"`
}
