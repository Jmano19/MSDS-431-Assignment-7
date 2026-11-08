package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Jmano19/MSDS-431-Assignment-7/crawler"
	"github.com/Jmano19/MSDS-431-Assignment-7/models"
	"github.com/Jmano19/MSDS-431-Assignment-7/scraper"
)

// Wikipedia URLs provided in the assignment.
var targetURLs = []string{
	"https://en.wikipedia.org/wiki/Robotics",
	"https://en.wikipedia.org/wiki/Robot",
	"https://en.wikipedia.org/wiki/Reinforcement_learning",
	"https://en.wikipedia.org/wiki/Robot_Operating_System",
	"https://en.wikipedia.org/wiki/Intelligent_agent",
	"https://en.wikipedia.org/wiki/Software_agent",
	"https://en.wikipedia.org/wiki/Robotic_process_automation",
	"https://en.wikipedia.org/wiki/Chatbot",
	"https://en.wikipedia.org/wiki/Applications_of_artificial_intelligence",
	"https://en.wikipedia.org/wiki/Android_(robot)",
}

const (
	outputPath  = "output/wikipedia_articles.jl"
	workerCount = 3
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run executes the crawling, scraping, and JSON Lines output pipeline.
func run() error {
	startTime := time.Now()

	fmt.Println("Starting Wikipedia crawler...")
	fmt.Printf("Target pages: %d\n", len(targetURLs))

	// Step 1: Download webpages concurrently using Colly.
	fmt.Println("\nStep 1: Crawling webpages...")

	pages, crawlErr := crawler.Crawl(targetURLs)

	if crawlErr != nil {
		log.Printf("Some crawling requests failed: %v", crawlErr)
	}

	fmt.Printf("Successfully downloaded: %d/%d pages\n",
		len(pages), len(targetURLs))

	if len(pages) == 0 {
		return fmt.Errorf("no webpages were successfully downloaded")
	}

	// Step 2: Extract readable article text concurrently.
	fmt.Println("\nStep 2: Extracting article text...")

	articles, scrapeFailures := scrapePages(pages, workerCount)

	for _, article := range articles {
		fmt.Printf("Scraped: %s\n", article.Title)
	}

	fmt.Printf("Successfully scraped: %d/%d downloaded pages\n",
		len(articles), len(pages))

	if len(articles) == 0 {
		return fmt.Errorf("no articles were successfully scraped")
	}

	// Step 3: Create the output directory if it does not exist.
	fmt.Println("\nStep 3: Writing JSON Lines output...")

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	// Create or overwrite the JSON Lines output file.
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer outputFile.Close()

	// Encode each article as one JSON object per line.
	encoder := json.NewEncoder(outputFile)

	for _, article := range articles {
		if err := encoder.Encode(article); err != nil {
			return fmt.Errorf("write article %q: %w",
				article.Title, err)
		}
	}

	// Step 4: Report execution metrics.
	elapsed := time.Since(startTime)

	fmt.Println("\nCrawl completed.")
	fmt.Println("-----------------------------------")
	fmt.Printf("Target pages:       %d\n", len(targetURLs))
	fmt.Printf("Downloaded pages:   %d\n", len(pages))
	fmt.Printf("Scraped articles:   %d\n", len(articles))
	fmt.Printf("Crawl failures:     %d\n", len(targetURLs)-len(pages))
	fmt.Printf("Scrape failures:    %d\n", scrapeFailures)
	fmt.Printf("Total failures:     %d\n",
		len(targetURLs)-len(pages)+scrapeFailures)
	fmt.Printf("Execution time:     %s\n", elapsed)
	fmt.Printf("Output file:        %s\n", outputPath)
	fmt.Println("-----------------------------------")

	return nil
}

// scrapeResult stores the outcome of processing one webpage.
type scrapeResult struct {
	article models.Article
	err     error
}

// scrapePages extracts article text using a concurrent worker pool.
func scrapePages(pages []crawler.Page, workers int) ([]models.Article, int) {
	if len(pages) == 0 {
		return nil, 0
	}

	if workers < 1 {
		workers = 1
	}

	if workers > len(pages) {
		workers = len(pages)
	}

	jobs := make(chan int)
	results := make([]scrapeResult, len(pages))

	var wg sync.WaitGroup

	// Start the scraping workers.
	for worker := 0; worker < workers; worker++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for index := range jobs {
				article, err := scraper.Extract(
					pages[index].URL,
					pages[index].HTML,
				)

				results[index] = scrapeResult{
					article: article,
					err:     err,
				}
			}
		}()
	}

	// Distribute webpage indexes to the workers.
	for index := range pages {
		jobs <- index
	}

	close(jobs)
	wg.Wait()

	// Collect successful results in their original order.
	var articles []models.Article
	var failures int

	for _, result := range results {
		if result.err != nil {
			log.Printf("Scraping failed: %v", result.err)
			failures++
			continue
		}

		articles = append(articles, result.article)
	}

	return articles, failures
}
