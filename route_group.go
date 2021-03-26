package group

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Middleware is just type alias of middleware function.
//
// Middleware は単にミドルウェア関数の型の別名.
type Middleware func(httprouter.Handle) httprouter.Handle

// RouteGroup is struct to retain group information.
//
// RouteGroup はグループ情報を保持する構造体.
type RouteGroup struct {
	children    []*RouteGroup
	handlers    map[string]httprouter.Handle
	middlewares []Middleware
	path        string
}

// New returns a new initialized RouteGroup.
// Define only base path first, other info define by method chaining.
//
// New は初期化された RouteGroup を返す.
// 初めはパスのみ定義し, 他の情報はメソッドチェーンで定義していく.
func New(path string) *RouteGroup {
	return &RouteGroup{
		children:    make([]*RouteGroup, 0),
		handlers:    make(map[string]httprouter.Handle),
		middlewares: make([]Middleware, 0),
		path:        path,
	}
}

// DELETE is a shortcut for (*RouteGroup).Handle(http.MethodDelete, handle).
func (r *RouteGroup) DELETE(handle httprouter.Handle) *RouteGroup {
	return r.Handle(http.MethodDelete, handle)
}

// GET is a shortcut for (*RouteGroup).Handle(http.MethodGet, handle).
func (r *RouteGroup) GET(handle httprouter.Handle) *RouteGroup {
	return r.Handle(http.MethodGet, handle)
}

// HEAD is a shortcut for (*RouteGroup).Handle(http.MethodHead, handle).
func (r *RouteGroup) HEAD(handle httprouter.Handle) *RouteGroup {
	return r.Handle(http.MethodHead, handle)
}

// OPTIONS is a shortcut for (*RouteGroup).Handle(http.MethodOptions, handle).
func (r *RouteGroup) OPTIONS(handle httprouter.Handle) *RouteGroup {
	return r.Handle(http.MethodOptions, handle)
}

// PATCH is a shortcut for (*RouteGroup).Handle(http.MethodPatch, handle).
func (r *RouteGroup) PATCH(handle httprouter.Handle) *RouteGroup {
	return r.Handle(http.MethodPatch, handle)
}

// POST is a shortcut for (*RouteGroup).Handle(http.MethodPost, handle).
func (r *RouteGroup) POST(handle httprouter.Handle) *RouteGroup {
	return r.Handle(http.MethodPost, handle)
}

// PUT is a shortcut for (*RouteGroup).Handle(http.MethodPut, handle).
func (r *RouteGroup) PUT(handle httprouter.Handle) *RouteGroup {
	return r.Handle(http.MethodPut, handle)
}

// Children returns self includes specified groups.
//
// Children は指定されたグループを含む自身を返す.
func (r *RouteGroup) Children(children ...*RouteGroup) *RouteGroup {
	r.children = append(r.children, children...)

	return r
}

// Handle returns self includes specified pair of HTTP method and handler.
// Handler will be overwritten if a registered HTTP method is specified.
//
// Handle は指定された HTTP メソッドとハンドラのペアを含む自身を返す.
// 登録済みの HTTP メソッドが指定された場合, ハンドラは上書きされる.
func (r *RouteGroup) Handle(method string, handler httprouter.Handle) *RouteGroup {
	r.handlers[method] = handler

	return r
}

// Middleware returns self includes specified middlewares.
//
// Middleware は指定されたミドルウェアを含む自身を返す.
func (r *RouteGroup) Middleware(middlewares ...Middleware) *RouteGroup {
	r.middlewares = append(r.middlewares, middlewares...)

	return r
}

// List returns slice of the pair of registered HTTP method
// and URI sorted by string. This includes the information of children.
//
// List は登録済みの HTTP メソッドと URI のペアを文字列でソートした slice を返す.
// 自身を起点とした子階層の情報も取得して返す.
//
// Example.
//
//     []string{
//         "GET     /",
//         "GET     /users",
//         "DELETE  /users/:id",
//         "GET     /users/:id",
//         "PUT     /users/:id",
//     }
//
func (r *RouteGroup) List() []string {
	// 8 = utf8.RuneCountInString(http.MethodOptions) + 1.
	const format = "%-8s%s"

	routes := r.list("")
	list := make([]string, 0, len(routes))

	sort.Sort(_sortBy(routes))

	for i := range routes {
		list = append(list, fmt.Sprintf(format, routes[i][0], routes[i][1]))
	}

	return list
}

func (r *RouteGroup) list(parentPath string) [][2]string {
	path := joinPath(parentPath, r.path)
	routes := make([][2]string, 0, r.len())

	for m := range r.handlers {
		routes = append(routes, [2]string{m, path})
	}

	for _, rg := range r.children {
		routes = append(routes, rg.list(path)...)
	}

	return routes
}

func (r *RouteGroup) len() int {
	cnt := len(r.handlers)

	for _, rg := range r.children {
		cnt += rg.len()
	}

	return cnt
}

type _sortBy [][2]string

func (a _sortBy) Len() int      { return len(a) }
func (a _sortBy) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a _sortBy) Less(i, j int) bool {
	if a[i][1] == a[j][1] {
		return a[i][0] < a[j][0]
	}

	return a[i][1] < a[j][1]
}

// Register set registered handlers to Router recursively.
// Middlewares are registered so that it is executed in
// the order of registration of the parent hierarchy,
// followed by the order of registration of the child hierarchies.
//
// Register は登録済みのハンドラを再帰的に Router に登録する.
// ミドルウェアは親階層の登録順, 続いて子階層の登録順で実行されるように登録される.
func (r *RouteGroup) Register(router *httprouter.Router) {
	r.registerChild(router, "", nil)
}

func (r *RouteGroup) registerChild(
	router *httprouter.Router,
	parentPath string,
	parentMiddlewares []Middleware,
) {
	ms := append(parentMiddlewares, r.middlewares...)
	path := joinPath(parentPath, r.path)

	for m, h := range r.handlers {
		router.Handle(m, path, middlewareWith(h, ms...))
	}

	for _, rg := range r.children {
		rg.registerChild(router, path, ms)
	}
}

func joinPath(ps ...string) string {
	var buf strings.Builder

	for _, p := range ps {
		if p := strings.Trim(p, "/"); p != "" {
			buf.WriteByte('/')
			buf.WriteString(p)
		}
	}

	if 0 < buf.Len() {
		return buf.String()
	}

	return "/"
}

func middlewareWith(h httprouter.Handle, ms ...Middleware) httprouter.Handle {
	for i := len(ms) - 1; 0 <= i; i-- {
		h = ms[i](h)
	}

	return h
}
