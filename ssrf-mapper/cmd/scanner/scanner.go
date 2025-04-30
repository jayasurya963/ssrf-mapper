package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

func testURL(targetURL string) {
	resp, err := http.Get(targetURL)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	body := string(bodyBytes)

	// Check for JavaScript-based redirects
	jsRedirect := regexp.MustCompile(`(?i)window\.location(?:\.href)?\s*=\s*["']([^"']+)["']`)
	if jsRedirect.MatchString(body) {
		fmt.Printf("[!] Client-side redirect detected in response: %s\n", targetURL)
		matches := jsRedirect.FindStringSubmatch(body)
		if len(matches) > 1 {
			fmt.Printf(" → Redirects to: %s\n", matches[1])
		}
		return
	}

	fmt.Printf("[OK] No client-side redirect found: %s\n", targetURL)
}
