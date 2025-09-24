package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ScrapeResult struct {
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Headers     []string `json:"headers"`
}

// fetchWithRetry tries to GET the URL with timeout + retry logic
func fetchWithRetry(url string) (*http.Response, error) {
	client := &http.Client{
		Timeout: 10 * time.Second, // request timeout
	}

	var resp *http.Response
	var err error
	for i := 1; i <= 3; i++ {
		resp, err = client.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		// log failed attempt
		log.Printf("Attempt %d failed for %s: %v", i, url, err)
		if resp != nil {
			resp.Body.Close()
		}

		time.Sleep(2 * time.Second) // wait before retry
	}
	return nil, fmt.Errorf("failed after 3 retries: %v", err)
}

func scrapePage(url string, wg *sync.WaitGroup, ch chan<- ScrapeResult) {
	defer wg.Done()

	resp, err := fetchWithRetry(url)
	if err != nil {
		log.Printf("Skipping %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Failed to parse %s: %v", url, err)
		return
	}

	title := strings.TrimSpace(doc.Find("title").Text())
	description, _ := doc.Find("meta[name='description']").Attr("content")

	var headers []string
	doc.Find("h1").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			headers = append(headers, text)
		}
	})

	ch <- ScrapeResult{
		URL:         url,
		Title:       title,
		Description: description,
		Headers:     headers,
	}
}

func writeResultsToJSONFile(results []ScrapeResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

func main() {
	// --- Ask user for input URLs ---
	fmt.Println("Enter URLs to scrape (comma-separated):")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		log.Fatal("No URLs provided, exiting.")
	}
	urls := strings.Split(input, ",")

	// --- Ask for output filename ---
	fmt.Print("Enter output filename (default: scraped_results.json): ")
	outputReader := bufio.NewReader(os.Stdin)
	outputInput, _ := outputReader.ReadString('\n')
	outputInput = strings.TrimSpace(outputInput)
	if outputInput == "" {
		outputInput = "scraped_results.json"
	}

	var wg sync.WaitGroup
	ch := make(chan ScrapeResult)
	var results []ScrapeResult
	var mu sync.Mutex

	// collect results
	go func() {
		for result := range ch {
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}
	}()

	// launch scrapers
	for _, url := range urls {
		wg.Add(1)
		go scrapePage(strings.TrimSpace(url), &wg, ch)
	}

	wg.Wait()
	close(ch)

	// write to file
	err := writeResultsToJSONFile(results, outputInput)
	if err != nil {
		log.Fatalf("Failed to write results to file: %v", err)
	}
	fmt.Printf("✅ Scraping completed. Results written to %s\n", outputInput)
}
