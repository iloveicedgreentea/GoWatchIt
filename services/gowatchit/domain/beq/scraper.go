package beq

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/iloveicedgreentea/gowatchit/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/net/html"
)

const catalogueURL = "https://beqcatalogue.readthedocs.io/en/latest/"

// AuthorScraper scrapes BEQ author names from the catalogue website
type AuthorScraper struct {
	httpClient *http.Client
}

// NewAuthorScraper creates a new AuthorScraper instance
func NewAuthorScraper() *AuthorScraper {
	return &AuthorScraper{
		httpClient: &http.Client{},
	}
}

// ScrapeAuthors fetches and parses the BEQ catalogue to extract author names
func (s *AuthorScraper) ScrapeAuthors(ctx context.Context) ([]string, error) {
	log := logger.GetLoggerFromContext(ctx)
	log.Debug("Fetching BEQ catalogue", zap.String("url", catalogueURL))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, catalogueURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch catalogue: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error("Failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	authors, err := s.parseAuthors(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse authors: %w", err)
	}

	log.Info("Successfully scraped authors", zap.Int("count", len(authors)), zap.Strings("authors", authors))
	return authors, nil
}

// parseAuthors extracts author names from the HTML content
func (s *AuthorScraper) parseAuthors(htmlContent string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	authors := make([]string, 0)
	seenAuthors := make(map[string]bool)

	// Recursively search for author links
	var findAuthors func(*html.Node)
	findAuthors = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			// Check if this is an author link by looking at the href attribute
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					// Author links have patterns like "aron7awol/", "mobe1969/", etc.
					// They are relative links ending with /
					href := attr.Val
					if strings.HasSuffix(href, "/") && !strings.Contains(href, "http") && !strings.Contains(href, "..") {
						authorName := strings.TrimSuffix(href, "/")
						// Filter out common navigation links
						if authorName != "" && !isNavigationLink(authorName) {
							if !seenAuthors[authorName] {
								authors = append(authors, authorName)
								seenAuthors[authorName] = true
							}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findAuthors(c)
		}
	}

	findAuthors(doc)

	if len(authors) == 0 {
		return nil, fmt.Errorf("no authors found in catalogue")
	}

	return authors, nil
}

// isNavigationLink filters out common navigation/structural links
func isNavigationLink(link string) bool {
	navigationLinks := []string{
		"index",
		"search",
		"genindex",
		"catalogue",
		"#",
	}

	for _, nav := range navigationLinks {
		if link == nav || strings.HasPrefix(link, nav) {
			return true
		}
	}

	return false
}
