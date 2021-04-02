package group

import (
	"fmt"
	"sort"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Route is struct to retain routing information.
// Route does not keep the information of
// the parent hierarchy (path, middleware).
//
// Route はルーティングの情報を保持する構造体.
// 親階層の情報（パス, ミドルウェア）は保持しない.
type Route struct {
	handler httprouter.Handle
	method  string
	path    string
}

// Handler just returns handler.
//
// Handler は単にハンドラを返す.
func (r *Route) Handler() httprouter.Handle { return r.handler }

// Method just returns HTTP method.
//
// Method は単に HTTP メソッドを返す.
func (r *Route) Method() string { return r.method }

// Path just returns path.
//
// Path は単にパスを返す.
func (r *Route) Path() string { return r.path }

// Routes is aggregation of `Route`.
//
// Routes は `Route` の集合体.
type Routes []*Route

// String sorts the registered HTTP method and URI pairs by string
// and returns them as a newline-separated string.
// This includes the information of children.
//
// String は登録済みの HTTP メソッドと URI のペアを
// 文字列でソートし, 改行区切りの文字列で返す.
// 自身を起点とした子階層の情報も取得して返す.
//
// Example.
//
//     `GET     /
//     GET     /users
//     DELETE  /users/:id
//     GET     /users/:id
//     PUT     /users/:id`
//
func (r Routes) String() string {
	// 8 = utf8.RuneCountInString(http.MethodOptions) + 1.
	const format = "%-8s%s"

	list := make([]string, 0, len(r))

	sort.Sort(r)

	for i := range r {
		list = append(list, fmt.Sprintf(format, r[i].method, r[i].path))
	}

	return strings.Join(list, "\n")
}

// Len is an implementation of `sort.Interface`.
//
// Len は `sort.Interface` の実装.
func (r Routes) Len() int { return len(r) }

// Swap is an implementation of `sort.Interface`.
//
// Swap は `sort.Interface` の実装.
func (r Routes) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

// Less is an implementation of `sort.Interface`.
// Sort by path, HTTP method priority.
//
// Less は `sort.Interface` の実装.
// パス, HTTP メソッドの優先順でソートする.
func (r Routes) Less(i, j int) bool {
	if r[i].path == r[j].path {
		return r[i].method < r[j].method
	}

	return r[i].path < r[j].path
}
