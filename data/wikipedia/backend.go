package wikipedia

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"snippets/data"

	"github.com/hashicorp/go-hclog"
)

// WikipediaBackend is a backend implementation that uses the wikipedia API.
type WikipediaBackend struct {
	l          hclog.Logger
	apiBaseURL string
}

// NewBackend creates a new WikipediaBackend.
func NewBackend(l hclog.Logger, apiBaseURL string) (*WikipediaBackend, error) {
	wb := &WikipediaBackend{
		l:          l,
		apiBaseURL: apiBaseURL,
	}
	return wb, nil
}

// GetSnippet gets a snippet from the backend.
func (wb *WikipediaBackend) GetSnippet(opts *data.GetSnippetOptions) (string, error) {

	urlEncodedName := url.QueryEscape(opts.Name)

	// Fetch the data from the backend
	url := fmt.Sprintf("%s?action=query&prop=revisions&titles=%s&rvlimit=1&formatversion=2&format=json&rvprop=content", wb.apiBaseURL, urlEncodedName)
	resp, err := http.Get(url)
	if err != nil {
		wb.l.Error("Error connecting to backend", "error", err)
		return "", data.ErrCannotConnectBackend
	}
	defer resp.Body.Close()

	// check the status code
	if resp.StatusCode == http.StatusNotFound {
		return "", data.ErrPageNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return "", data.ErrUnexpectedBackendResponse
	}

	// Define the struct to hold only the parts of response that we are interested
	type BackendResponse struct {
		Query struct {
			Pages []struct {
				Revisions []struct {
					Content string `json:"content"`
				} `json:"revisions"`
				Missing bool `json:"missing"`
			} `json:"pages"`
		} `json:"query"`
	}

	var backendResponse BackendResponse

	// Decode the response body into a struct
	err = data.FromJSON(&backendResponse, resp.Body)
	if err != nil {
		return "", data.ErrUnexpectedBackendResponse
	}

	if len(backendResponse.Query.Pages) == 0 {
		return "", data.ErrPageNotFound
	}

	if backendResponse.Query.Pages[0].Missing {
		return "", data.ErrPageNotFound
	}

	if len(backendResponse.Query.Pages[0].Revisions) == 0 {
		return "", data.ErrPageNotFound
	}

	if len(backendResponse.Query.Pages[0].Revisions[0].Content) == 0 {
		return "", data.ErrPageNotFound
	}

	shortDescription, err := extractShortDescription(backendResponse.Query.Pages[0].Revisions[0].Content)
	if err != nil {
		return "", err
	}

	return shortDescription, nil

}

// extractShortDescription extracts the short description from the page wikitext content.
func extractShortDescription(text string) (string, error) {

	re := regexp.MustCompile(`{{Short description\|(.*?)}}`)
	match := re.FindStringSubmatch(text)
	if len(match) != 2 {
		return "", data.ErrShortDescriptionNotFound
	}
	return match[1], nil

}
