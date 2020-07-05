// No implementation.
// This package is only used in router file (default: "router.go").
// See example.

package stdrouter

import "net/http"


type Router struct{}

func NewRouter() Router { return Router{} }

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request)                         {}
func (router Router) HandleFunc(path interface{}, method interface{}, handlerFunc interface{}) {}
func (router Router) HandleNotFound(handlerFunc interface{})                                   {}
func (router Router) HandleMethodNotAllowed(handlerFunc interface{})                           {}
