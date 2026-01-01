//go:build ignore
// +build ignore

// This program scrapes license text from opensource.org and updates the plain text files.
// Run with: go generate ./internal/licenses/...
//
// Usage:
//
//	go run scrape_licenses.go [flags]
//	  -licenses string
//	        Path to licenses.json file (default "data/licenses.json")
//	  -output string
//	        Output directory for plain text files (default "data/texts/plain")
//	  -force
//	        Force re-download of all licenses (default: skip existing)
//	  -id string
//	        Scrape only the specified license ID
//	  -update-licenses
//	        Update licenses.json from upstream if changed (default: true)
//	  -verbose
//	        Print verbose output
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const licensesURL = "https://s3.amazonaws.com/api.opensource.org/licenses/licenses.json"

// License represents a license entry from licenses.json
type License struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	licensesPath := flag.String("licenses", "data/licenses.json", "Path to licenses.json file")
	outputDir := flag.String("output", "data/texts/plain", "Output directory for plain text files")
	force := flag.Bool("force", false, "Force re-download of all licenses")
	singleID := flag.String("id", "", "Scrape only the specified license ID")
	updateLicenses := flag.Bool("update-licenses", true, "Update licenses.json from upstream if changed")
	verbose := flag.Bool("verbose", false, "Print verbose output")
	flag.Parse()

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Update licenses.json if requested
	if *updateLicenses {
		updated, err := updateLicensesJSON(client, *licensesPath, *verbose)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to update licenses.json: %v\n", err)
		} else if updated {
			fmt.Println("Updated licenses.json from upstream")
		} else if *verbose {
			fmt.Println("licenses.json is up to date")
		}
	}

	// Load licenses.json
	licenses, err := loadLicenses(*licensesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading licenses: %v\n", err)
		os.Exit(1)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Filter to single ID if specified
	if *singleID != "" {
		filtered := make([]License, 0, 1)
		for _, lic := range licenses {
			if strings.EqualFold(lic.ID, *singleID) {
				filtered = append(filtered, lic)
				break
			}
		}
		if len(filtered) == 0 {
			fmt.Fprintf(os.Stderr, "License ID '%s' not found in licenses.json\n", *singleID)
			os.Exit(1)
		}
		licenses = filtered
	}

	// Process each license
	successCount := 0
	skipCount := 0
	errorCount := 0

	for _, lic := range licenses {
		outputPath := filepath.Join(*outputDir, lic.ID)

		// Skip if file exists and not forcing
		if !*force {
			if _, err := os.Stat(outputPath); err == nil {
				if *verbose {
					fmt.Printf("SKIP: %s (already exists)\n", lic.ID)
				}
				skipCount++
				continue
			}
		}

		if *verbose {
			fmt.Printf("Fetching: %s (%s)\n", lic.ID, lic.Name)
		}

		// Try both URL formats (new and old)
		text, err := scrapeLicense(client, lic.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s: %v\n", lic.ID, err)
			errorCount++
			continue
		}

		// Write the license text
		if err := os.WriteFile(outputPath, []byte(text), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s: failed to write file: %v\n", lic.ID, err)
			errorCount++
			continue
		}

		fmt.Printf("OK: %s\n", lic.ID)
		successCount++

		// Rate limiting to be nice to the server
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Printf("\nSummary: %d succeeded, %d skipped, %d errors\n", successCount, skipCount, errorCount)
	if errorCount > 0 {
		os.Exit(1)
	}
}

func loadLicenses(path string) ([]License, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var licenses []License
	if err := json.Unmarshal(data, &licenses); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	return licenses, nil
}

// updateLicensesJSON downloads licenses.json from upstream if it has changed.
// It uses ETag/Last-Modified headers to avoid re-downloading unchanged content.
// Returns true if the file was updated, false if unchanged.
func updateLicensesJSON(client *http.Client, destPath string, verbose bool) (bool, error) {
	etagPath := destPath + ".etag"

	// Build request with conditional headers
	req, err := http.NewRequest("GET", licensesURL, nil)
	if err != nil {
		return false, fmt.Errorf("creating request: %w", err)
	}

	// Add If-None-Match header if we have a cached ETag
	if etag, err := os.ReadFile(etagPath); err == nil {
		req.Header.Set("If-None-Match", strings.TrimSpace(string(etag)))
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// 304 Not Modified - file hasn't changed
	if resp.StatusCode == http.StatusNotModified {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("reading response: %w", err)
	}

	// Validate it's valid JSON
	var licenses []License
	if err := json.Unmarshal(body, &licenses); err != nil {
		return false, fmt.Errorf("invalid JSON from upstream: %w", err)
	}

	if verbose {
		fmt.Printf("Downloaded %d licenses from upstream\n", len(licenses))
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return false, fmt.Errorf("creating directory: %w", err)
	}

	// Write the file with pretty formatting
	prettyJSON, err := json.MarshalIndent(licenses, "", "  ")
	if err != nil {
		return false, fmt.Errorf("formatting JSON: %w", err)
	}

	if err := os.WriteFile(destPath, prettyJSON, 0644); err != nil {
		return false, fmt.Errorf("writing file: %w", err)
	}

	// Save ETag for future conditional requests
	if etag := resp.Header.Get("ETag"); etag != "" {
		if err := os.WriteFile(etagPath, []byte(etag), 0644); err != nil {
			// Non-fatal, just log if verbose
			if verbose {
				fmt.Printf("Warning: failed to save ETag: %v\n", err)
			}
		}
	}

	return true, nil
}

func scrapeLicense(client *http.Client, id string) (string, error) {
	// Try new URL format first: /license/mit
	urls := []string{
		fmt.Sprintf("https://opensource.org/license/%s", strings.ToLower(id)),
		fmt.Sprintf("https://opensource.org/licenses/%s", id),
	}

	var lastErr error
	for _, url := range urls {
		text, err := fetchAndExtract(client, url)
		if err == nil {
			return text, nil
		}
		lastErr = err
	}

	return "", fmt.Errorf("failed to fetch from all URLs: %w", lastErr)
}

func fetchAndExtract(client *http.Client, url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %w", err)
	}

	// Parse HTML and extract license text
	text, err := extractLicenseText(string(body))
	if err != nil {
		return "", err
	}

	return text, nil
}

func extractLicenseText(htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", fmt.Errorf("parsing HTML: %w", err)
	}

	// Find the license-content div
	licenseDiv := findElementByClass(doc, "license-content")
	if licenseDiv == nil {
		// Fallback to entry-content
		licenseDiv = findElementByClass(doc, "entry-content")
	}
	if licenseDiv == nil {
		return "", fmt.Errorf("could not find license content in page")
	}

	// Extract text content
	var textBuilder strings.Builder
	extractText(licenseDiv, &textBuilder)

	// Clean up the text
	text := cleanLicenseText(textBuilder.String())

	if len(strings.TrimSpace(text)) == 0 {
		return "", fmt.Errorf("extracted empty license text")
	}

	return text, nil
}

func findElementByClass(n *html.Node, class string) *html.Node {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, class) {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findElementByClass(c, class); result != nil {
			return result
		}
	}

	return nil
}

func extractText(n *html.Node, builder *strings.Builder) {
	if n.Type == html.TextNode {
		builder.WriteString(n.Data)
	}

	// Add newlines for block elements
	if n.Type == html.ElementNode {
		switch n.Data {
		case "p", "div", "br", "li", "h1", "h2", "h3", "h4", "h5", "h6":
			builder.WriteString("\n")
		}
	}

	// Skip sidebar and widget content
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "class" {
				if strings.Contains(attr.Val, "sidebar") ||
					strings.Contains(attr.Val, "widget") ||
					strings.Contains(attr.Val, "syndication") {
					return
				}
			}
			if attr.Key == "role" && attr.Val == "complementary" {
				return
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, builder)
	}

	if n.Type == html.ElementNode {
		switch n.Data {
		case "p", "div", "li":
			builder.WriteString("\n")
		}
	}
}

// wrapText wraps text at the specified column width, preferring to break at word boundaries.
// It preserves paragraph breaks (empty lines) and handles lines that are already shorter than the width.
func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	lines := strings.Split(text, "\n")
	var wrapped []string

	for _, line := range lines {
		// Preserve empty lines (paragraph breaks)
		if strings.TrimSpace(line) == "" {
			wrapped = append(wrapped, "")
			continue
		}

		// If line is already short enough, keep it as-is
		if len(line) <= width {
			wrapped = append(wrapped, line)
			continue
		}

		// Wrap long lines
		words := strings.Fields(line)
		if len(words) == 0 {
			wrapped = append(wrapped, "")
			continue
		}

		var currentLine strings.Builder
		currentLine.WriteString(words[0])

		for _, word := range words[1:] {
			// Check if adding this word would exceed the width
			// +1 for the space between words
			if currentLine.Len()+1+len(word) > width {
				// Flush the current line if it has content
				if currentLine.Len() > 0 {
					wrapped = append(wrapped, currentLine.String())
					currentLine.Reset()
				}

				// If word is longer than width, add it on its own line and continue
				if len(word) > width {
					wrapped = append(wrapped, word)
					continue
				}
				
				// Start new line with this word
				currentLine.WriteString(word)
			} else {
				// Add space and word to current line
				currentLine.WriteString(" ")
				currentLine.WriteString(word)
			}
		}

		// Don't forget the last line
		if currentLine.Len() > 0 {
			wrapped = append(wrapped, currentLine.String())
		}
	}

	return strings.Join(wrapped, "\n")
}

func cleanLicenseText(text string) string {
	// Decode HTML entities
	text = decodeHTMLEntities(text)

	// Normalize whitespace
	lines := strings.Split(text, "\n")
	var cleanedLines []string
	prevEmpty := false

	for _, line := range lines {
		// Trim whitespace from each line
		line = strings.TrimSpace(line)

		// Skip multiple consecutive empty lines
		if line == "" {
			if !prevEmpty {
				cleanedLines = append(cleanedLines, "")
				prevEmpty = true
			}
			continue
		}
		prevEmpty = false
		cleanedLines = append(cleanedLines, line)
	}

	// Join and trim
	result := strings.Join(cleanedLines, "\n")
	result = strings.TrimSpace(result)

	// Wrap long lines at a reasonable column width
	result = wrapText(result, 78)

	// Ensure trailing newline
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}

	return result
}

func decodeHTMLEntities(text string) string {
	replacements := map[string]string{
		"&#8220;": "\"",
		"&#8221;": "\"",
		"&#8216;": "'",
		"&#8217;": "'",
		"&#8211;": "-",
		"&#8212;": "--",
		"&quot;":  "\"",
		"&amp;":   "&",
		"&lt;":    "<",
		"&gt;":    ">",
		"&nbsp;":  " ",
		"&#160;":  " ",
		"&copy;":  "(c)",
		"&#169;":  "(c)",
	}

	for entity, replacement := range replacements {
		text = strings.ReplaceAll(text, entity, replacement)
	}

	// Handle numeric entities like &#8220;
	numericEntity := regexp.MustCompile(`&#(\d+);`)
	text = numericEntity.ReplaceAllStringFunc(text, func(match string) string {
		var num int
		fmt.Sscanf(match, "&#%d;", &num)
		if num > 0 && num < 128 {
			return string(rune(num))
		}
		// Return original for high unicode points we haven't mapped
		return match
	})

	return text
}
