# Concurrent Web Scraper in Go

A simple **concurrent web scraper** written in Go, capable of fetching multiple URLs in parallel, extracting the page title, description, and `<h1>` headers, and saving the results to a JSON file. Includes retry logic and request timeout handling.

---

## Features

- Concurrent scraping using Go routines and channels.
- Fetches:
  - Page **title**
  - Meta **description**
  - All `<h1>` headers
- Retry logic for failed HTTP requests (up to 3 attempts per URL).
- Customizable output filename.
- Saves results in **JSON format**.

---

## Requirements

- Go 1.20+ (or compatible)
- [goquery](https://github.com/PuerkitoBio/goquery) for HTML parsing

Install dependencies:

```bash
go get github.com/PuerkitoBio/goquery
```

## Usage

1. Clone the repository:

```bash
git clone https://github.com/your-username/concurrent-scraper.git
cd concurrent-scraper
```
2. Run the scraper

```bash
go run main.go
```
3. Enter the URLs to scrape (comma-separated), e.g.:
```bash
https://example.com,https://golang.org
```
4. Enter the output filename (or press enter to use the default)
5. The results will be save in .json file with the chosen name or the default if you didn't choose one :)



