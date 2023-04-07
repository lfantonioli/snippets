package handlers

import (
	"snippets/data"

	"github.com/hashicorp/go-hclog"
)

// SnippetsService is a service that handles snippets requests.
type SnippetsService struct {
	l       hclog.Logger
	backend data.SnippetsBackend
}

// NewSnippetsService creates a new SnippetsService.
func NewSnippetsService(l hclog.Logger, backend data.SnippetsBackend) *SnippetsService {
	return &SnippetsService{l, backend}
}
