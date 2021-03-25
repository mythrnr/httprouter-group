package group

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_joinPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args []string
		want string
	}{{
		args: []string{"/"},
		want: "/",
	}, {
		args: []string{""},
		want: "/",
	}, {
		args: []string{"///"},
		want: "/",
	}, {
		args: []string{"/", "/"},
		want: "/",
	}, {
		args: []string{"", ""},
		want: "/",
	}, {
		args: []string{"/", "/users"},
		want: "/users",
	}, {
		args: []string{"", "users"},
		want: "/users",
	}, {
		args: []string{"/users"},
		want: "/users",
	}, {
		args: []string{"users"},
		want: "/users",
	}, {
		args: []string{"users", "/"},
		want: "/users",
	}, {
		args: []string{"users", ""},
		want: "/users",
	}, {
		args: []string{"/", "/users", "/:id"},
		want: "/users/:id",
	}, {
		args: []string{"", "users", ":id"},
		want: "/users/:id",
	}, {
		args: []string{"/users", "/:id"},
		want: "/users/:id",
	}, {
		args: []string{"users", ":id"},
		want: "/users/:id",
	}, {
		args: []string{"users/:id"},
		want: "/users/:id",
	}, {
		args: []string{"users", "", "/:id"},
		want: "/users/:id",
	}}

	for _, tt := range tests {
		assert.Equal(t, tt.want, joinPath(tt.args...))
	}
}

func Test_middlewareWith(t *testing.T) {
	t.Parallel()

	ms := []Middleware{
		func(h httprouter.Handle) httprouter.Handle {
			return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
				rw.Write([]byte("1"))
				h(rw, r, p)
				rw.Write([]byte("1"))
			}
		},
		func(h httprouter.Handle) httprouter.Handle {
			return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
				rw.Write([]byte("2"))
				h(rw, r, p)
				rw.Write([]byte("2"))
			}
		},
		func(h httprouter.Handle) httprouter.Handle {
			return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
				rw.Write([]byte("3"))
				h(rw, r, p)
				rw.Write([]byte("3"))
			}
		},
		func(h httprouter.Handle) httprouter.Handle {
			return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
				rw.Write([]byte("4"))
				h(rw, r, p)
				rw.Write([]byte("4"))
			}
		},
	}

	h := func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		rw.Write([]byte("5"))
	}

	h = middlewareWith(h, ms...)

	rw := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/test",
		nil,
	)
	p := httprouter.Params{}

	h(rw, req, p)

	body, err := ioutil.ReadAll(rw.Body)
	require.Nil(t, err)

	assert.Equal(t, "123454321", string(body))
}
