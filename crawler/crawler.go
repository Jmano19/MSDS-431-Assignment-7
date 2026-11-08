package crawler

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

// Page represents a downloaded webpage before text extraction.
type Page struct {
	URL  string
	HTML []byte
}

// Crawl downloads the supplied URLs concurrently and returns their HTML.
func Crawl(urls []string) ([]Page, error) {
	if len(urls) == 0 {
		return nil, nil
	}

	// Initialize Colly with asynchronous requests.
	collector := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("MSDS431WikipediaCrawler/1.0"),
	)

	// Limit concurrent requests to avoid overwhelming Wikipedia.
	if err := collector.Limit(&colly.LimitRule{
		DomainGlob:  "en.wikipedia.org",
		Parallelism: 3,
		Delay:       500 * time.Millisecond,
	}); err != nil {
		return nil, fmt.Errorf("configure crawler: %w", err)
	}

	var (
		mu          sync.Mutex
		pages       []Page
		crawlErrors []error
	)

	// Capture successfully downloaded webpage content.
	collector.OnResponse(func(response *colly.Response) {
		page := Page{
			URL:  response.Request.URL.String(),
			HTML: append([]byte(nil), response.Body...),
		}

		mu.Lock()
		pages = append(pages, page)
		mu.Unlock()
	})

	// Record request errors without stopping the remaining crawl.
	collector.OnError(func(response *colly.Response, err error) {
		mu.Lock()
		defer mu.Unlock()

		crawlErrors = append(
			crawlErrors,
			fmt.Errorf("crawl %s: %w", response.Request.URL, err),
		)
	})

	// Preserve the original URL order in the final results.
	urlOrder := make(map[string]int, len(urls))
	for index, targetURL := range urls {
		if _, exists := urlOrder[targetURL]; !exists {
			urlOrder[targetURL] = index
		}

		if err := collector.Visit(targetURL); err != nil {
			mu.Lock()
			crawlErrors = append(
				crawlErrors,
				fmt.Errorf("visit %s: %w", targetURL, err),
			)
			mu.Unlock()
		}
	}

	// Wait until all asynchronous requests finish.
	collector.Wait()

	// Sort successfully downloaded pages into their original order.
	sort.SliceStable(pages, func(i, j int) bool {
		return urlOrder[pages[i].URL] < urlOrder[pages[j].URL]
	})

	return pages, errors.Join(crawlErrors...)
}
