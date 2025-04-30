package main

import (
    "bufio"
    "encoding/json"
    "flag"
    "fmt"
    "net/http"
    "os"
    "strings"
    "sync"
    "time"
	"log"
	"ssrf-mapper/scanner"  //Import your scanner package for ScanURL

    "github.com/yourusername/ssrf-mapper/internal/utils"
)

type Result struct {
    Payload      string `json:"payload"`
    StatusCode   int    `json:"status_code"`
    Cloud        string `json:"cloud"`
    RedirectedTo string `json:"redirected_to,omitempty"`
}

func main() {
    target := flag.String("target", "", "Target URL with INJECT_HERE")
    wordlist := flag.String("wordlist", "", "Wordlist of SSRF payloads")
    header := flag.String("header", "", "Custom headers (semicolon-separated)")
    concurrency := flag.Int("concurrency", 10, "Number of concurrent workers")
    output := flag.String("output", "results.json", "Output file (JSON)")

    flag.Parse()

    if *target == "" || *wordlist == "" {
        fmt.Println("Usage: ssrf-mapper --target <url> --wordlist <file>")
        os.Exit(1)
    }

    targets := prepareTargets(*wordlist)
    jobs := make(chan string, len(targets))
    results := make(chan Result, len(targets))

    var wg sync.WaitGroup
    for i := 0; i < *concurrency; i++ {
        wg.Add(1)
        go worker(jobs, results, *target, *header, &wg)
    }

    for _, ip := range targets {
        encoded := utils.GenerateEncodedPayloads(ip)
        for _, enc := range encoded {
            jobs <- enc
        }
    }
    close(jobs)

    wg.Wait()
    close(results)

    saveResults(results, *output)
}

func prepareTargets(wordlistPath string) []string {
    file, err := os.Open(wordlistPath)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, strings.TrimSpace(scanner.Text()))
    }
    return lines
}

func worker(jobs <-chan string, results chan<- Result, targetTemplate, header string, wg *sync.WaitGroup) {
    defer wg.Done()
    client := &http.Client{
        Timeout:       5 * time.Second,
        CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
    }

    headers := parseHeaders(header)

    for payload := range jobs {
        url := strings.Replace(targetTemplate, "INJECT_HERE", payload, -1)
        req, _ := http.NewRequest("GET", url, nil)
        for k, v := range headers {
            req.Header.Set(k, v)
        }

        resp, err := client.Do(req)
        if err != nil {
            continue
        }
        defer resp.Body.Close()

        result := Result{
            Payload:    payload,
            StatusCode: resp.StatusCode,
            Cloud:      utils.DetectCloud(resp),
        }

        if resp.StatusCode >= 300 && resp.StatusCode < 400 {
            loc, err := resp.Location()
            if err == nil {
                result.RedirectedTo = loc.String()
            }
        }

        results <- result
    }
}

func parseHeaders(raw string) map[string]string {
    headers := make(map[string]string)
    if raw == "" {
        return headers
    }
    pairs := strings.Split(raw, ";")
    for _, pair := range pairs {
        kv := strings.SplitN(pair, ":", 2)
        if len(kv) == 2 {
            headers[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
        }
    }
    return headers
}

func saveResults(results <-chan Result, path string) {
    file, err := os.Create(path)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    json.NewEncoder(file).Encode(resultsToSlice(results))
}

func resultsToSlice(results <-chan Result) []Result {
    var all []Result
    for res := range results {
        all = append(all, res)
    }
    return all
}

func main() {
	// Define flags for command-line arguments
	var detectJSRedirects bool
	var targetURL string
	
	// Flag to enable/disable JavaScript redirect detection
	flag.BoolVar(&detectJSRedirects, "detect-js-redirects", true, "Enable client-side redirect detection")
	
	// Flag to pass the target URL you want to scan
	flag.StringVar(&targetURL, "url", "", "Target URL to test")
	
	// Parse the flags
	flag.Parse()

	// Ensure the target URL is provided
	if targetURL == "" {
		log.Fatal("[ERROR] Please provide a target URL using the -url flag.")
	}

	// Call the ScanURL function from the scanner package
	scanner.ScanURL(targetURL, detectJSRedirects)

	// Print scan completion info
	fmt.Printf("[INFO] Scan complete for: %s\n", targetURL)
}