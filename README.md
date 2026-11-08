# MSDS 431 – Assignment 7: Concurrent Web Crawler

## Project Overview

This project implements a concurrent web crawler and web scraper in Go using the&#x20;

[Colly](https://github.com/gocolly/colly) framework.

The application downloads and extracts text from ten Wikipedia articles related to robotics, artificial intelligence, intelligent agents, and automation. The extracted content is cleaned, structured, and exported into a JSON Lines (`.jl`) file.

The project demonstrates Go's concurrency features, including goroutines, channels, synchronization, and concurrent processing.

## Objectives

- Build a web crawler using the Colly framework.
- Download multiple webpages concurrently.
- Extract readable article text while removing unnecessary HTML elements.
- Process downloaded webpages using a concurrent worker pool.
- Export structured article data in JSON Lines format.
- Test application functionality and measure performance.
- Compare sequential and concurrent processing.

## Technologies Used

| Technology         | Purpose                           |
| ------------------ | --------------------------------- |
| Go                 | Primary programming language      |
| Colly              | Concurrent web crawling           |
| goquery            | HTML parsing and text extraction  |
| Goroutines         | Concurrent article processing     |
| Channels           | Distributing work between workers |
| sync.WaitGroup     | Coordinating concurrent workers   |
| sync.Mutex         | Protecting shared crawler results |
| encoding/json      | JSON Lines serialization          |
| Go testing package | Unit tests and benchmarks         |
| httptest           | Local HTTP testing                |

## Project Structure

```text
MSDS-431-Assignment-7/
│
├── crawler/
│   ├── crawler.go
│   └── crawler_test.go
│
├── models/
│   └── article.go
│
├── scraper/
│   ├── scraper.go
│   └── scraper_test.go
│
├── output/
│   └── wikipedia_articles.jl
│
├── benchmarks/
│   └── benchmark_results.txt
│
├── main.go
├── main_test.go
├── go.mod
├── go.sum
├── README.md
├── GENAI_CONVERSATION_LOG.md
└── wikipedia-crawler.exe
```

## Application Architecture

The application consists of four primary components.

### 1. Crawler

The `crawler` package uses Colly to download the target webpages.

- Asynchronous requests are enabled using `colly.Async(true)`.
- Requests are limited to three concurrent connections.
- A 500-millisecond delay is configured between requests.
- A mutex protects shared webpage results and error collections.
- Downloaded pages are sorted to preserve the original URL order.

### 2. Scraper

The `scraper` package uses goquery to parse downloaded HTML.

It extracts:

- Article titles.
- Article paragraphs.
- Section headings.
- Relevant list items.

It removes unnecessary elements, including navigation, scripts, tables, infoboxes, and other metadata.

### 3. Models

The `models` package defines the `Article` structure used throughout the application.

Each article contains:

- `title`: Wikipedia article title.
- `url`: Original webpage URL.
- `text`: Extracted article text.
- `scraped_at`: UTC timestamp indicating when the article was processed.

### 4. Main Application

The main application coordinates crawling, extraction, and output generation.

A worker pool of three goroutines processes downloaded webpages concurrently. A channel distributes webpage indexes to the workers, and a `sync.WaitGroup` ensures all workers finish before results are collected.

Successful results are written to the output file in their original order.

## Target Websites

The application crawls the following ten Wikipedia articles:

1. [Robotics](https://en.wikipedia.org/wiki/Robotics)
2. [Robot](https://en.wikipedia.org/wiki/Robot)
3. [Reinforcement Learning](https://en.wikipedia.org/wiki/Reinforcement_learning)
4. [Robot Operating System](https://en.wikipedia.org/wiki/Robot_Operating_System)
5. [Intelligent Agent](https://en.wikipedia.org/wiki/Intelligent_agent)
6. [Software Agent](https://en.wikipedia.org/wiki/Software_agent)
7. [Robotic Process Automation](https://en.wikipedia.org/wiki/Robotic_process_automation)
8. [Chatbot](https://en.wikipedia.org/wiki/Chatbot)
9. [Applications of Artificial Intelligence](https://en.wikipedia.org/wiki/Applications_of_artificial_intelligence)
10. [Android (robot)](https://en.wikipedia.org/wiki/Android_\(robot\))

## Installation and Execution

### Requirements

- Go 1.27 or compatible version.
- Internet connection for downloading Wikipedia webpages.

### 1. Clone the repository

```powershell
git clone https://github.com/Jmano19/MSDS-431-Assignment-7.git
cd MSDS-431-Assignment-7
```

### 2. Install dependencies

```powershell
go mod tidy
```

### 3. Run the application

```powershell
go run .
```

The application downloads the webpages, extracts article text, and creates the JSON Lines output file.

### 4. Build the Windows executable

```powershell
go build -o wikipedia-crawler.exe .
```

Run the executable:

```powershell
.\wikipedia-crawler.exe
```

## Output Format

The application generates:

```text
output/wikipedia_articles.jl
```

The output uses JSON Lines format, meaning each line contains one independent JSON object representing a scraped article.

Example:

```json
{"title":"Robotics","url":"https://en.wikipedia.org/wiki/Robotics","text":"Robotics is an interdisciplinary field...","scraped_at":"2026-11-07T..."}
```

The output file can be processed using standard JSON Lines tools or imported into downstream data processing workflows.

## Testing

Unit tests are included for the crawler and scraper packages, as well as the concurrent scraping worker pool.

Run all tests:

```powershell
go test -v ./...
```

Run tests with coverage:

```powershell
go test -cover ./...
```

Run static analysis:

```powershell
go vet ./...
```

The crawler tests use Go's `httptest` package to test HTTP requests locally without relying on external websites.

Test coverage includes:

- Successful webpage downloads.
- Empty URL input.
- Invalid URLs.
- Successful HTML extraction.
- Missing article titles.
- Missing article content.
- Empty HTML content.
- Exclusion of unwanted HTML elements.
- Concurrent article extraction.
- Error handling and preservation of result order.

The race detector was not executed successfully in the development environment because CGO was disabled and a compatible C compiler was unavailable.

## Performance Benchmarks

The project includes benchmarks comparing sequential and concurrent article extraction.

Run benchmarks:

```powershell
go test -bench=BenchmarkScrape -benchmem -count=5 .
```

### Benchmark Environment

| Component              | Specification                      |
| ---------------------- | ---------------------------------- |
| Operating System       | Windows                            |
| Architecture           | AMD64                              |
| Processor              | AMD Ryzen 7 5800X 8-Core Processor |
| Go Version             | 1.27.1                             |
| Articles per iteration | 10                                 |
| Concurrent workers     | 3                                  |

### Benchmark Results

Results below represent the average of five benchmark trials.

| Metric           |    Sequential |    Concurrent |
| ---------------- | ------------: | ------------: |
| Execution time   | 456,948 ns/op | 267,130 ns/op |
| Memory allocated |  398,885 B/op |  402,624 B/op |
| Allocations      |         2,840 |         2,851 |

The raw results from all five benchmark trials are available in `benchmarks/benchmark_results.txt`.

### Performance Analysis

Concurrent article extraction achieved approximately **1.71x the processing speed** of sequential extraction.

This represents an approximately **41.5% reduction in execution time**.

The concurrent implementation used approximately 0.94% more memory and 11 additional allocations per benchmark iteration.

These results demonstrate the performance benefits of concurrent processing when handling multiple independent webpages.

The benchmarks use synthetic HTML documents to isolate and compare extraction performance. They do not measure network download performance.

## Actual Crawl Results

The application was successfully executed against all ten assigned Wikipedia URLs.

| Metric                  |       Result |
| ----------------------- | -----------: |
| Target webpages         |           10 |
| Successfully downloaded |           10 |
| Successfully scraped    |           10 |
| Crawl failures          |            0 |
| Scraping failures       |            0 |
| JSON Lines records      |           10 |
| Total execution time    | 2.31 seconds |

All ten articles were successfully extracted and written to `output/wikipedia_articles.jl`.

## GenAI Tools and Development

Generative AI tools were used as development assistants during the implementation of this project.

AI assistance supported:

- Go application structure and package organization.
- Colly crawler implementation.
- HTML extraction and parsing.
- Concurrent worker-pool implementation.
- Unit testing and benchmark development.
- Debugging and documentation.

Generated code was reviewed, executed, and tested during development. The developer remains responsible for the final implementation and verification.

A separate `GENAI_CONVERSATION_LOG.md` file documents the use of generative AI tools during development.

## Conclusion

This project demonstrates how Go can be used to build a concurrent web crawler and scraper.

Colly provides an effective framework for asynchronous webpage retrieval, while goroutines and channels enable concurrent processing of downloaded content.

The benchmark results demonstrate that concurrent extraction improved processing speed by approximately 41.5% compared with sequential extraction, with a relatively small increase in memory usage.

The final application successfully crawled all ten assigned Wikipedia pages, extracted their article text, and generated a structured JSON Lines output file.
