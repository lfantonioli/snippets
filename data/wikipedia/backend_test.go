package wikipedia

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"snippets/data"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func createNewWikipediaTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.URL.Query().Get("titles") {
		case "John_Carmacky":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
{
  "batchcomplete": true,
  "query": {
	"normalized": [
	  {
		"fromencoded": false,
		"from": "John_Carmacky",
		"to": "John Carmacky"
	  }
	],
	"pages": [
	  {
		"ns": 0,
		"title": "John Carmacky",
		"missing": true
	  }
	]
  }
}`))
		case "John_Carmack", "John Carmack":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
{
  "continue": {
	"rvcontinue": "20230324134545|1146374185",
	"continue": "||"
  },
  "warnings": {
	"main": {
	  "warnings": "Subscribe to the mediawiki-api-announce mailing list at <https://lists.wikimedia.org/postorius/lists/mediawiki-api-announce.lists.wikimedia.org/> for notice of API deprecations and breaking changes. Use [[Special:ApiFeatureUsage]] to see usage of deprecated features by your application."
	},
	"revisions": {
	  "warnings": "Because \"rvslots\" was not specified, a legacy format has been used for the output. This format is deprecated, and in the future the new format will always be used."
	}
  },
  "query": {
	"normalized": [
	  {
		"fromencoded": false,
		"from": "John_Carmack",
		"to": "John Carmack"
	  }
	],
	"pages": [
	  {
		"pageid": 38368,
		"ns": 0,
		"title": "John Carmack",
		"revisions": [
		  {
			"contentformat": "text/x-wiki",
			"contentmodel": "wikitext",
			"content": "{{Short description|American computer programmer and video game developer}}\n{{Other people|John Carmack}}\n{{Use mdy dates|date=September 2019}}\n{{Infobox person\n| name               = John Carmack\n| image              = John Carmack GDC 2010.jpg\n| caption            = Carmack in 2010\n| birth_date         = {{Birth date and age|1970|8|21}}<ref name=\"wired.com\" />\n| birth_place i"
		  }
		]
	  }
	]
  }
}
			`))
		case "Yoshua_Bengio", "Yoshua Bengio":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
{
  "continue": {
	"rvcontinue": "20230319105228|1145481799",
	"continue": "||"
  },
  "query": {
	"normalized": [
	  {
		"fromencoded": false,
		"from": "Yoshua_Bengio",
		"to": "Yoshua Bengio"
	  }
	],
	"pages": [
	  {
		"pageid": 47749536,
		"ns": 0,
		"title": "Yoshua Bengio",
		"revisions": [
		  {
			"contentformat": "text/x-wiki",
			"contentmodel": "wikitext",
			"content": "Hey--{{Short description|Canadian computer scientist}}\n{{Use mdy dates|date=March 2019}}\n{{Infobox scientist\n|name              = Yoshua Bengio\n|image             = Yoshua Bengio 2019 cropped.jpg\n|image_size        = \n| honorific_suffix = {{post-nominals|country=CAN|OC|FRS|FRSC|size=100}}\n|caption           = Yoshua Bengio in 2019\n|birth_date"
		  }
		]
	  }
	]
  }
}
			`))

		case "John_No_Short_description":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
{
  "continue": {
	"rvcontinue": "20230319105228|1145481799",
	"continue": "||"
  },
  "query": {
	"normalized": [
	  {
		"fromencoded": false,
		"from": "Yoshua_Bengio",
		"to": "Yoshua Bengio"
	  }
	],
	"pages": [
	  {
		"pageid": 47749536,
		"ns": 0,
		"title": "Yoshua Bengio",
		"revisions": [
		  {
			"contentformat": "text/x-wiki",
			"contentmodel": "wikitext",
			"content": "Hey--{{Canadian computer scientist}}\n{{Use mdy dates|date=March 2019}}\n{{Infobox scientist\n|name              = Yoshua Bengio\n|image             = Yoshua Bengio 2019 cropped.jpg\n|image_size        = \n| honorific_suffix = {{post-nominals|country=CAN|OC|FRS|FRSC|size=100}}\n|caption           = Yoshua Bengio in 2019\n|birth_date"
		  }
		]
	  }
	]
  }
}
			`))

		case "John_Unexpected_Response":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
{
  "quereY": {
	"pages": [
	  {
		"pageid": 47749536,
			`))

		default:
			w.WriteHeader(http.StatusNotFound)
		}

	}))
}

func TestWikipediaBackend_GetSnippet(t *testing.T) {
	// Set up a test server that will always return the same response
	testServer := createNewWikipediaTestServer()
	defer testServer.Close()

	backend, err := NewBackend(hclog.NewNullLogger(), testServer.URL)
	assert.NoError(t, err)

	tests := []struct {
		name             string
		options          *data.GetSnippetOptions
		expectedResponse string
		expectedError    error
	}{
		{
			name: "Yoshua_Bengio valid page test",
			options: &data.GetSnippetOptions{
				Name: "Yoshua_Bengio",
			},
			expectedResponse: "Canadian computer scientist",
			expectedError:    nil,
		},
		{
			name: "John Carmack valid page test that uses space",
			options: &data.GetSnippetOptions{
				Name: "John Carmack",
			},
			expectedResponse: "American computer programmer and video game developer",
			expectedError:    nil,
		},
		{
			name: "invalid page test",
			options: &data.GetSnippetOptions{
				Name: "Not_Yosua_Bengio",
			},
			expectedResponse: "",
			expectedError:    data.ErrPageNotFound,
		},
		{
			name: "invalid page test John_Carmacky",
			options: &data.GetSnippetOptions{
				Name: "John_Carmacky",
			},
			expectedResponse: "",
			expectedError:    data.ErrPageNotFound,
		},
		{
			name: "no short description page test",
			options: &data.GetSnippetOptions{
				Name: "John_No_Short_description",
			},
			expectedResponse: "",
			expectedError:    data.ErrShortDescriptionNotFound,
		},
		{
			name: "unexpected backend response",
			options: &data.GetSnippetOptions{
				Name: "John_Unexpected_Response",
			},
			expectedResponse: "",
			expectedError:    data.ErrUnexpectedBackendResponse,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := backend.GetSnippet(test.options)
			assert.Equal(t, test.expectedResponse, response)
			assert.Equal(t, test.expectedError, err)
		})
	}
}
