package handlers

import (
	"net/http"
	"net/url"
	"snippets/data"

	"github.com/gorilla/mux"
)

// GetSnippetResponse represents the response from the GetSnippetHandler
type GetSnippetResponse struct {
	ShortDescription string `json:"short_description"`
	ErrorMsg         string `json:"error_message"`
}

// GetSnippetHandler returns a handler function that handles GET requests to /snippets/{name}
func (ss *SnippetsService) GetSnippetHandler(rw http.ResponseWriter, r *http.Request) {
	ss.l.Debug("Get snippet")

	// get the name from the URL
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		ss.l.Error("Error getting name from URL")
		writeErrorResponse(rw, "Missing Person Name", http.StatusBadRequest)
		return
	}

	// unescape the name
	unescapedName, err := url.QueryUnescape(name)
	if err != nil {
		ss.l.Error("Error unescaping name", "error", err)
		writeErrorResponse(rw, "Error unescaping name", http.StatusBadRequest)
		return
	}

	opts := &data.GetSnippetOptions{
		Name: unescapedName,
	}

	// get the snippet from the backend
	snippet, err := ss.backend.GetSnippet(opts)
	if err != nil {
		switch err {
		case data.ErrPageNotFound:
			ss.l.Info("Could not fulfil request: Page Not Found", "query", name, "error", err)
			writeErrorResponse(rw, "Page Not Found", http.StatusNotFound)
		case data.ErrShortDescriptionNotFound:
			ss.l.Info("Could not fulfil request: Short Description Not Found", "query", name, "error", err)
			writeErrorResponse(rw, "Short Description Not Found", http.StatusUnprocessableEntity)
		case data.ErrUnexpectedBackendResponse:
			ss.l.Info("Could not fulfil request: Unexpected Backend Response", "query", name, "error", err)
			writeErrorResponse(rw, "Internal Server Error", http.StatusInternalServerError)
		case data.ErrCannotConnectBackend:
			ss.l.Info("Could not fulfil request: Cannot Connect to Backend. Check API URL", "query", name, "error", err)
			writeErrorResponse(rw, "Internal Server Error", http.StatusInternalServerError)
		default:
			ss.l.Info("Could not fulfil request: Generic Error on retrieval of Short Description", "query", name, "error", err)
			writeErrorResponse(rw, "Error retrieving Short Description", http.StatusInternalServerError)
		}

		return
	}

	response := &GetSnippetResponse{
		ShortDescription: snippet,
	}

	// serialize the response and write it to the response writer
	rw.Header().Add("Content-Type", "application/json")
	err = data.ToJSON(response, rw)
	if err != nil {
		ss.l.Error("Error serializing snippet response", "error", err)
		http.Error(rw, "{}", http.StatusInternalServerError)
		return
	}

}

// writeErrorResponse writes an error response to the response writer
func writeErrorResponse(rw http.ResponseWriter, errMsg string, code int) {
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(code)

	// we ignore the error that possibly can come from data.ToJSON here because there is nothing we can do about it
	_ = data.ToJSON(
		&GetSnippetResponse{
			ErrorMsg: errMsg,
		},
		rw,
	)
}
