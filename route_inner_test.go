package group

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sortBy(t *testing.T) {
	t.Parallel()

	expected := strings.Join([]string{
		"GET     /",
		"OPTIONS /",
		"PUT     /",
		"GET     /users",
		"DELETE  /users/:id",
		"GET     /users/:id",
	}, "\n")

	actual := Routes([]*Route{
		{method: http.MethodDelete, path: "/users/:id", handler: nil},
		{method: http.MethodGet, path: "/", handler: nil},
		{method: http.MethodGet, path: "/users", handler: nil},
		{method: http.MethodOptions, path: "/", handler: nil},
		{method: http.MethodGet, path: "/users/:id", handler: nil},
		{method: http.MethodPut, path: "/", handler: nil},
	})

	assert.Equal(t, expected, actual.String())
}
