package group

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sortBy(t *testing.T) {
	t.Parallel()

	expected := []*Route{
		{method: "GET", path: "/", handler: nil},
		{method: "PUT", path: "/", handler: nil},
		{method: "GET", path: "/users", handler: nil},
		{method: "DELETE", path: "/users/:id", handler: nil},
		{method: "GET", path: "/users/:id", handler: nil},
	}

	actual := Routes([]*Route{
		{method: "DELETE", path: "/users/:id", handler: nil},
		{method: "GET", path: "/", handler: nil},
		{method: "GET", path: "/users", handler: nil},
		{method: "GET", path: "/users/:id", handler: nil},
		{method: "PUT", path: "/", handler: nil},
	})

	sort.Sort(actual)

	assert.Equal(t, expected, []*Route(actual))
}
