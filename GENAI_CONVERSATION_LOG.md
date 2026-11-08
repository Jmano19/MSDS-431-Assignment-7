# GenAI Conversation Log – MSDS 431 Assignment 7

## Overview

OpenAI ChatGPT was used as a generative AI development assistant during the implementation of this Go web crawler and scraper.

The purpose of using GenAI was to support software design, code development, debugging, testing, benchmarking, and documentation. Generated code and recommendations were reviewed and tested locally before being incorporated into the project.

The developer made the implementation decisions, executed the code, evaluated test results, and verified the final output.

## GenAI Tool

- **Tool:** OpenAI ChatGPT
- **Model:** GPT-6
- **Primary use:** Go programming assistance, debugging, testing, performance analysis, and technical documentation.

## Development Interaction Log

### 1. Project Setup and Architecture

**Objective:** Establish the Go project structure and dependencies.

**GenAI assistance:**
- Recommended organizing the project into `crawler`, `scraper`, `models`, `output`, and `benchmarks` directories.
- Provided guidance for initializing the Go module.
- Recommended Colly for web crawling and goquery for HTML parsing.
- Explained the roles of the different application packages.

**Developer actions:**
- Initialized the Go module.
- Installed Colly and its dependencies.
- Created the project directories and initial source files.

### 2. Data Model Implementation

**Objective:** Define a consistent structure for scraped article data.

**GenAI assistance:**
- Suggested an `Article` structure containing the article title, URL, extracted text, and scrape timestamp.
- Explained the use of JSON struct tags for serialization.

**Developer actions:**
- Created `models/article.go`.
- Formatted the source code and updated dependencies.

### 3. Web Crawler Implementation

**Objective:** Download the ten assigned Wikipedia webpages.

**GenAI assistance:**
- Generated a Colly-based crawler implementation.
- Explained asynchronous requests and request concurrency limits.
- Recommended synchronization using a mutex to protect shared results.
- Suggested preserving the original URL order and collecting errors without terminating the entire crawl.

**Developer actions:**
- Implemented `crawler/crawler.go`.
- Configured asynchronous crawling with a concurrency limit of three and a 500-millisecond delay.
- Executed the crawler against the assigned Wikipedia URLs.

**Verification:** All ten webpages were successfully downloaded.

### 4. HTML Scraper Implementation

**Objective:** Extract readable article content and remove unnecessary HTML elements.

**GenAI assistance:**
- Provided goquery-based HTML extraction logic.
- Recommended extracting article titles, paragraphs, headings, and list items.
- Suggested removing navigation, tables, infoboxes, scripts, and metadata.
- Recommended returning descriptive errors for missing or invalid article content.

**Developer actions:**
- Implemented `scraper/scraper.go`.
- Created unit tests covering successful extraction and invalid HTML scenarios.

**Verification:** All scraper tests passed, with 97.1% statement coverage.

### 5. Concurrent Processing and JSON Lines Output

**Objective:** Integrate crawling, extraction, and structured output generation.

**GenAI assistance:**
- Suggested a worker-pool design using goroutines and channels.
- Explained the use of `sync.WaitGroup` to coordinate workers.
- Provided guidance for preserving result order while processing webpages concurrently.
- Recommended using `encoding/json` to write one JSON object per line.

**Developer actions:**
- Implemented the application workflow in `main.go`.
- Configured three concurrent scraping workers.
- Generated `output/wikipedia_articles.jl`.

**Verification:** All ten articles were successfully scraped and written to the output file. The JSON Lines file contained ten valid records.

### 6. Crawler Testing

**Objective:** Test crawler behavior without relying on external websites.

**GenAI assistance:**
- Recommended using Go's `httptest` package.
- Suggested tests for successful downloads, empty input, and invalid URLs.

**Developer actions:**
- Created `crawler/crawler_test.go`.
- Executed the crawler tests and measured statement coverage.

**Verification:** All crawler tests passed, with 82.5% statement coverage.

### 7. Concurrent Worker Testing and Benchmarking

**Objective:** Verify concurrent processing and compare sequential versus concurrent performance.

**GenAI assistance:**
- Suggested testing the worker pool with both valid and invalid HTML.
- Provided benchmark implementations for sequential and concurrent article extraction.
- Explained Go benchmark metrics, including execution time, memory allocation, and allocation counts.
- Assisted with troubleshooting benchmark discovery and execution.

**Developer actions:**
- Created `main_test.go`.
- Added `TestScrapePages`.
- Implemented `BenchmarkScrapeSequential` and `BenchmarkScrapeConcurrent`.
- Executed five benchmark trials for each implementation.

**Verification:** Both benchmarks passed.

Average benchmark results:

| Metric | Sequential | Concurrent |
|---|---:|---:|
| Execution time | 460,298 ns/op | 255,586 ns/op |
| Memory allocated | 398,885 B/op | 402,626 B/op |
| Allocations | 2,840 | 2,851 |

Concurrent processing was approximately 1.80 times faster, representing a 44.5% reduction in execution time.

### 8. Debugging and Code Quality

**Objective:** Verify the application and identify potential issues.

**GenAI assistance:**
- Helped troubleshoot benchmark execution.
- Recommended running Go's static analysis tools.
- Explained the CGO requirement for Go's race detector.

**Developer actions:**
- Ran `go test ./...`.
- Ran `go test -cover ./...`.
- Ran `go vet ./...`.
- Attempted `go test -race ./...`.

**Verification:** Unit tests and static analysis passed. The race detector could not execute because CGO was disabled and a compatible C compiler was unavailable in the development environment.

### 9. Documentation

**Objective:** Document the project implementation, architecture, testing, and performance.

**GenAI assistance:**
- Helped structure the README.
- Assisted with documenting dependencies, application execution, output format, test coverage, and benchmark results.
- Helped summarize the performance comparison.

**Developer actions:**
- Created the project README.
- Prepared this GenAI conversation log.

## Final Developer Verification

The developer executed and verified the application locally.

Final observed application results:

- Ten target webpages.
- Ten successful downloads.
- Ten successfully extracted articles.
- Zero crawl failures.
- Zero scraping failures.
- Ten JSON Lines records.
- Approximately 3.25 seconds total execution time.

The project tests and static analysis also completed successfully.

## Reflection

GenAI was useful for accelerating development and understanding Go-specific concepts, particularly Colly's asynchronous crawling, goroutine-based worker pools, channels, synchronization, and benchmarking.

An important part of the development process was validating generated code through actual execution, unit tests, benchmark results, and inspection of the output file.

The project demonstrates how GenAI can assist software engineering while keeping implementation decisions, testing, and verification under developer control.