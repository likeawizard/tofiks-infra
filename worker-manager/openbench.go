package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HasActiveWork checks the OpenBench index page for active tests.
// It looks for the "Active Tests" section and checks if it contains test entries.
func HasActiveWork(openbenchURL string) (bool, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(openbenchURL)
	if err != nil {
		return false, fmt.Errorf("fetching OpenBench index: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("OpenBench returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("reading OpenBench response: %w", err)
	}

	return detectActiveTests(string(body)), nil
}

// detectActiveTests parses the HTML to find active test entries.
// The OpenBench index page has a table header like "Active : 2 Machines / 7 Threads"
// followed by test rows, then "Finished" or "Pending" or "Awaiting" sections.
func detectActiveTests(html string) bool {
	// The active section header starts with "Active :" followed by machine/thread counts
	activeIdx := strings.Index(html, "Active :")
	if activeIdx == -1 {
		return false
	}

	// Get the content after "Active :" until the next section
	rest := html[activeIdx:]

	// The section ends at Finished, Pending, or Awaiting headers
	for _, marker := range []string{"Finished", "Pending", "Awaiting"} {
		if idx := strings.Index(rest, marker); idx > 0 {
			rest = rest[:idx]
			break
		}
	}

	// Active tests have links to /test/ pages
	return strings.Contains(rest, "/test/")
}
