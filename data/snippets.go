package data

import "errors"

var (
	ErrShortDescriptionNotFound  = errors.New("Short Description not found in page")
	ErrPageNotFound              = errors.New("Page not found")
	ErrCannotConnectBackend      = errors.New("Cannot connect to backend")
	ErrUnexpectedBackendResponse = errors.New("Unexpected backend response")
)

type GetSnippetOptions struct {
	Name string
}

type SnippetsBackend interface {
	GetSnippet(opts *GetSnippetOptions) (string, error)
}
