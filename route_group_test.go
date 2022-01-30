package group_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	group "github.com/mythrnr/httprouter-group"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RouteGroup(t *testing.T) {
	t.Parallel()

	g := group.New("/").Handle(
		http.MethodGet,
		func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
			rw.Write([]byte("GET /\n"))
		},
	).Handle(
		http.MethodPost,
		func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
			rw.Write([]byte("POST /\n"))
		},
	).Handle(
		http.MethodOptions,
		func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
			rw.Write([]byte("OPTIONS /\n"))
		},
	).Middleware(
		func(h httprouter.Handle) httprouter.Handle {
			return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
				rw.Write([]byte("Middleware 1: before\n"))
				h(rw, r, p)
				rw.Write([]byte("Middleware 1: after\n"))
			}
		},
		func(h httprouter.Handle) httprouter.Handle {
			return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
				rw.Write([]byte("Middleware 2: before\n"))
				h(rw, r, p)
				rw.Write([]byte("Middleware 2: after\n"))
			}
		},
	).Children(
		group.New("/users").Handle(
			http.MethodGet,
			func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
				rw.Write([]byte("GET /users\n"))
			},
		).Handle(
			http.MethodPost,
			func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
				rw.Write([]byte("POST /users\n"))
			},
		).Middleware(
			func(h httprouter.Handle) httprouter.Handle {
				return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
					rw.Write([]byte("Middleware 3: before\n"))
					h(rw, r, p)
					rw.Write([]byte("Middleware 3: after\n"))
				}
			},
			func(h httprouter.Handle) httprouter.Handle {
				return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
					rw.Write([]byte("Middleware 4: before\n"))
					h(rw, r, p)
					rw.Write([]byte("Middleware 4: after\n"))
				}
			},
		).Children(
			group.New("/:id").Handle(
				http.MethodGet,
				func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
					rw.Write([]byte("GET /users/:id\n"))
				},
			).Middleware(
				func(h httprouter.Handle) httprouter.Handle {
					return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
						rw.Write([]byte("Middleware 5: before\n"))
						h(rw, r, p)
						rw.Write([]byte("Middleware 5: after\n"))
					}
				},
				func(h httprouter.Handle) httprouter.Handle {
					return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
						rw.Write([]byte("Middleware 6: before\n"))
						h(rw, r, p)
						rw.Write([]byte("Middleware 6: after\n"))
					}
				},
			),
			group.New("/:id").Handle(
				http.MethodPut,
				func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
					rw.Write([]byte("PUT /users/:id\n"))
				},
			).Middleware(
				func(h httprouter.Handle) httprouter.Handle {
					return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
						rw.Write([]byte("Middleware 7: before\n"))
						h(rw, r, p)
						rw.Write([]byte("Middleware 7: after\n"))
					}
				},
				func(h httprouter.Handle) httprouter.Handle {
					return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
						rw.Write([]byte("Middleware 8: before\n"))
						h(rw, r, p)
						rw.Write([]byte("Middleware 8: after\n"))
					}
				},
			),
		),
	)

	assert.Equal(t, strings.Join([]string{
		"GET     /",
		"OPTIONS /",
		"POST    /",
		"GET     /users",
		"POST    /users",
		"GET     /users/:id",
		"PUT     /users/:id",
	}, "\n"), g.Routes().String())

	router := httprouter.New()

	for _, r := range g.Routes() {
		router.Handle(r.Method(), r.Path(), r.Handler())
	}

	{
		w := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodGet,
			"/",
			io.NopCloser(bytes.NewReader(nil)),
		)

		router.ServeHTTP(w, r)

		body, err := io.ReadAll(w.Body)
		require.Nil(t, err)

		values := strings.Split(string(body[:len(body)-1]), "\n")
		assert.Equal(t, []string{
			"Middleware 1: before",
			"Middleware 2: before",
			"GET /",
			"Middleware 2: after",
			"Middleware 1: after",
		}, values)
	}

	{
		w := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			"/",
			io.NopCloser(bytes.NewReader([]byte("body"))),
		)

		router.ServeHTTP(w, r)

		body, err := io.ReadAll(w.Body)
		require.Nil(t, err)

		values := strings.Split(string(body[:len(body)-1]), "\n")
		assert.Equal(t, []string{
			"Middleware 1: before",
			"Middleware 2: before",
			"POST /",
			"Middleware 2: after",
			"Middleware 1: after",
		}, values)
	}

	{
		w := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodGet,
			"/users",
			io.NopCloser(bytes.NewReader(nil)),
		)

		router.ServeHTTP(w, r)

		body, err := io.ReadAll(w.Body)
		require.Nil(t, err)

		values := strings.Split(string(body[:len(body)-1]), "\n")
		assert.Equal(t, []string{
			"Middleware 1: before",
			"Middleware 2: before",
			"Middleware 3: before",
			"Middleware 4: before",
			"GET /users",
			"Middleware 4: after",
			"Middleware 3: after",
			"Middleware 2: after",
			"Middleware 1: after",
		}, values)
	}

	{
		w := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			"/users",
			io.NopCloser(bytes.NewReader([]byte("body"))),
		)

		router.ServeHTTP(w, r)

		body, err := io.ReadAll(w.Body)
		require.Nil(t, err)

		values := strings.Split(string(body[:len(body)-1]), "\n")
		assert.Equal(t, []string{
			"Middleware 1: before",
			"Middleware 2: before",
			"Middleware 3: before",
			"Middleware 4: before",
			"POST /users",
			"Middleware 4: after",
			"Middleware 3: after",
			"Middleware 2: after",
			"Middleware 1: after",
		}, values)
	}

	{
		w := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodGet,
			"/users/1",
			io.NopCloser(bytes.NewReader(nil)),
		)

		router.ServeHTTP(w, r)

		body, err := io.ReadAll(w.Body)
		require.Nil(t, err)

		values := strings.Split(string(body[:len(body)-1]), "\n")
		assert.Equal(t, []string{
			"Middleware 1: before",
			"Middleware 2: before",
			"Middleware 3: before",
			"Middleware 4: before",
			"Middleware 5: before",
			"Middleware 6: before",
			"GET /users/:id",
			"Middleware 6: after",
			"Middleware 5: after",
			"Middleware 4: after",
			"Middleware 3: after",
			"Middleware 2: after",
			"Middleware 1: after",
		}, values)
	}

	{
		w := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodPut,
			"/users/1",
			io.NopCloser(bytes.NewReader([]byte("body"))),
		)

		router.ServeHTTP(w, r)

		body, err := io.ReadAll(w.Body)
		require.Nil(t, err)

		values := strings.Split(string(body[:len(body)-1]), "\n")
		assert.Equal(t, []string{
			"Middleware 1: before",
			"Middleware 2: before",
			"Middleware 3: before",
			"Middleware 4: before",
			"Middleware 7: before",
			"Middleware 8: before",
			"PUT /users/:id",
			"Middleware 8: after",
			"Middleware 7: after",
			"Middleware 4: after",
			"Middleware 3: after",
			"Middleware 2: after",
			"Middleware 1: after",
		}, values)
	}
}

func Test_RouteGroup_shortcut(t *testing.T) {
	t.Parallel()

	g := group.New(
		"/users/:id",
	).DELETE(
		func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			w.Write([]byte("DELETE /users/" + p.ByName("id")))
		},
	).GET(
		func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			w.Write([]byte("GET /users/" + p.ByName("id")))
		},
	).HEAD(
		func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			w.Write([]byte("HEAD /users/" + p.ByName("id")))
		},
	).OPTIONS(
		func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			w.Write([]byte("OPTIONS /users/" + p.ByName("id")))
		},
	).PATCH(
		func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			w.Write([]byte("PATCH /users/" + p.ByName("id")))
		},
	).POST(
		func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			w.Write([]byte("POST /users/" + p.ByName("id")))
		},
	).PUT(
		func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			w.Write([]byte("PUT /users/" + p.ByName("id")))
		},
	)

	assert.Equal(t, strings.Join([]string{
		"DELETE  /users/:id",
		"GET     /users/:id",
		"HEAD    /users/:id",
		"OPTIONS /users/:id",
		"PATCH   /users/:id",
		"POST    /users/:id",
		"PUT     /users/:id",
	}, "\n"), g.Routes().String())

	router := httprouter.New()

	for _, r := range g.Routes() {
		router.Handle(r.Method(), r.Path(), r.Handler())
	}

	{
		w := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodDelete,
			"/users/1",
			io.NopCloser(bytes.NewReader(nil)),
		)

		router.ServeHTTP(w, r)

		body, err := io.ReadAll(w.Body)

		require.Nil(t, err)
		assert.Equal(t, "DELETE /users/1", string(body))
	}

	{
		w := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodGet,
			"/users/2",
			io.NopCloser(bytes.NewReader(nil)),
		)

		router.ServeHTTP(w, r)

		body, err := io.ReadAll(w.Body)

		require.Nil(t, err)
		assert.Equal(t, "GET /users/2", string(body))
	}

	{
		w := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodHead,
			"/users/3",
			io.NopCloser(bytes.NewReader(nil)),
		)

		router.ServeHTTP(w, r)

		body, err := io.ReadAll(w.Body)

		require.Nil(t, err)
		assert.Equal(t, "HEAD /users/3", string(body))
	}

	{
		w := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodOptions,
			"/users/4",
			io.NopCloser(bytes.NewReader(nil)),
		)

		router.ServeHTTP(w, r)

		body, err := io.ReadAll(w.Body)

		require.Nil(t, err)
		assert.Equal(t, "OPTIONS /users/4", string(body))
	}

	{
		w := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodPatch,
			"/users/5",
			io.NopCloser(bytes.NewReader(nil)),
		)

		router.ServeHTTP(w, r)

		body, err := io.ReadAll(w.Body)

		require.Nil(t, err)
		assert.Equal(t, "PATCH /users/5", string(body))
	}

	{
		w := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			"/users/6",
			io.NopCloser(bytes.NewReader(nil)),
		)

		router.ServeHTTP(w, r)

		body, err := io.ReadAll(w.Body)

		require.Nil(t, err)
		assert.Equal(t, "POST /users/6", string(body))
	}

	{
		w := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodPut,
			"/users/7",
			io.NopCloser(bytes.NewReader(nil)),
		)

		router.ServeHTTP(w, r)

		body, err := io.ReadAll(w.Body)

		require.Nil(t, err)
		assert.Equal(t, "PUT /users/7", string(body))
	}
}

func Test_RouteGroup_Any(t *testing.T) {
	t.Parallel()

	g := group.New("/:param").GET(
		func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			w.Write([]byte("GET /" + p.ByName("param")))
		},
	).Any(
		func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			w.Write([]byte("Any: " + r.Method + " /" + p.ByName("param")))
		},
	).DELETE(
		func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			w.Write([]byte("DELETE /" + p.ByName("param")))
		},
	)

	assert.Equal(t, strings.Join([]string{
		"CONNECT /:param",
		"DELETE  /:param",
		"GET     /:param",
		"HEAD    /:param",
		"OPTIONS /:param",
		"PATCH   /:param",
		"POST    /:param",
		"PUT     /:param",
		"TRACE   /:param",
	}, "\n"), g.Routes().String())

	router := httprouter.New()

	for _, r := range g.Routes() {
		router.Handle(r.Method(), r.Path(), r.Handler())
	}

	for _, m := range []string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
		http.MethodTrace,
	} {
		w := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(
			context.Background(),
			m,
			"/parameter",
			io.NopCloser(bytes.NewReader(nil)),
		)

		router.ServeHTTP(w, r)

		body, err := io.ReadAll(w.Body)

		require.Nil(t, err)

		switch m {
		case http.MethodDelete, http.MethodGet:
			assert.Equal(t, m+" /parameter", string(body))
		default:
			assert.Equal(t, "Any: "+m+" /parameter", string(body))
		}
	}
}
