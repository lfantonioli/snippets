Snippets Service API
====================

The Snippets Service API provides a way to retrieve short descriptions of Wikipedia pages.

The service has three endpoints:

Get /snippets/{name}
Get /docs
Get /openapi.yaml

The service relies on a backend service (wikipedia API) to retrieve the short descriptions of Wikipedia pages.

Endpoints
---------

### Get Snippet
The /snippets/{name} endpoint is used to retrieve the short description of a Wikipedia page. The name parameter is the name of the page, URL-encoded.

Request:
  The request must be a GET request to the /snippets/{name} endpoint. The name parameter must be provided as a URL-encoded string.

Response:
  The response is a JSON object with the following fields:
    short_description (string): The short description of the Wikipedia page.
    error_message (string): If there was an error retrieving the short description, this field will contain an error message.
  If the requested page is not found or does not have a short description, an appropriate http status code will be returned.
    * 400 Bad Request: The request is malformed or missing required parameters.
    * 404 Not Found: The requested page could not be found or the endpoint does not exist.
    * 422 Unprocessable Entity: The requested page does not have a short description.
    * 500 Internal Server Error: There was an error retrieving the short description from the backend.

Example request:
  GET /snippets/John_Carmack


Example response:
  HTTP/1.1 200 OK
  Content-Type: application/json

  {
    "short_description": "American computer programmer and video game developer",
    "error_message": ""
  }

### Docs

The /docs endpoint provides an interactive web page powered by Redoc that allows to see the specification of the service

### OpenAPI spec

The  /openapi.yaml endpoint provides a simple OpenAPI Yaml specification file that documents the service



Instalation
-----------

### Option 1: Go development environment

If you have the Go development environment installed on your machine, you can compile the program using the following commands:

$ cd /path/to/snippets
$ go build .

This will generate a binary called snippets in the current directory. You can then run the program using the following command:

$ ./snippets

Alternatively, you can install the program using the following command:

$ go install .

This will install the binary in your Go binary directory (usually $GOPATH/bin or $HOME/go/bin). You can then run the program using the following command:

$ snippets 

### Option 2: Docker Image

If you don't have the Go development environment installed on your machine or prefer to use Docker, you can build a Docker image of the program using the provided Dockerfile. First, make sure Docker is installed on your machine. Then, run the following command from the root folder of the project:

$ docker build -t snippets .

This will build a Docker image called snippets. You can then run a Docker container with the program using the following command:

$ docker run -it --rm -p 9095:9095 snippets 

This will start a Docker container with the program running on port 9095 (it's default port).
