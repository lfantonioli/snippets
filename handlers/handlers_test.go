package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"

	"snippets/data"
)

func TestGetSnippetHandler(t *testing.T) {
	logger := hclog.NewNullLogger()
	backend := &MockBackend{}

	s := NewSnippetsService(logger, backend)

	cases := []struct {
		name     string
		query    string
		status   int
		response string
	}{
		{
			name:     "valid",
			query:    "Test",
			status:   http.StatusOK,
			response: "{\"short_description\":\"Test Snippet\",\"error_message\":\"\"}\n",
		},
		{
			name:     "not found",
			query:    "NotFound",
			status:   http.StatusNotFound,
			response: "{\"short_description\":\"\",\"error_message\":\"Page Not Found\"}\n",
		},
		{
			name:     "short description not found",
			query:    "NoShortDescription",
			status:   http.StatusUnprocessableEntity,
			response: "{\"short_description\":\"\",\"error_message\":\"Short Description Not Found\"}\n",
		},
		{
			name:     "unexpected backend response",
			query:    "UnexpectedResponse",
			status:   http.StatusInternalServerError,
			response: "{\"short_description\":\"\",\"error_message\":\"Internal Server Error\"}\n",
		},
		{
			name:     "cannot connect backend",
			query:    "CannotConnect",
			status:   http.StatusInternalServerError,
			response: "{\"short_description\":\"\",\"error_message\":\"Internal Server Error\"}\n",
		},
		{
			name:     "generic error",
			query:    "GenericError",
			status:   http.StatusInternalServerError,
			response: "{\"short_description\":\"\",\"error_message\":\"Error retrieving Short Description\"}\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/snippets/"+c.query, nil)
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/snippets/{name}", s.GetSnippetHandler).Methods("GET")

			backend.snippet = c.query

			router.ServeHTTP(rr, req)

			assert.Equal(t, c.status, rr.Code)

			expectedResponse := c.response

			p := rr.Body.String()

			assert.Equal(t, expectedResponse, p)
		})
	}
}

type MockBackend struct {
	snippet string
}

func (mb *MockBackend) GetSnippet(opts *data.GetSnippetOptions) (string, error) {
	switch mb.snippet {
	case "Test":
		return "Test Snippet", nil
	case "NoShortDescription":
		return "", data.ErrShortDescriptionNotFound
	case "UnexpectedResponse":
		return "", data.ErrUnexpectedBackendResponse
	case "CannotConnect":
		return "", data.ErrCannotConnectBackend
	case "GenericError":
		return "", errors.New("some error")
	default:
		return "", data.ErrPageNotFound
	}
}
