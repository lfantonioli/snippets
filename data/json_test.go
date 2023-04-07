package data

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToJSON(t *testing.T) {
	type Person struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}
	p := Person{Name: "John Doe", Age: 42, Email: "john.doe@example.com"}

	var buf bytes.Buffer
	err := ToJSON(p, &buf)
	assert.NoError(t, err)

	expected := `{"name":"John Doe","age":42,"email":"john.doe@example.com"}`
	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(buf.String()))
}

func TestFromJSON(t *testing.T) {
	type Person struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}
	jsonStr := `{"name":"John Doe","age":42,"email":"john.doe@example.com"}`

	var p Person
	err := FromJSON(&p, strings.NewReader(jsonStr))
	assert.NoError(t, err)

	assert.Equal(t, "John Doe", p.Name)
	assert.Equal(t, 42, p.Age)
	assert.Equal(t, "john.doe@example.com", p.Email)
}

func TestFromJSON_Error(t *testing.T) {
	type Person struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}
	jsonStr := `{"name":"John Doe","age":42,"email":}`

	var p Person
	err := FromJSON(&p, strings.NewReader(jsonStr))
	assert.Error(t, err)
}
