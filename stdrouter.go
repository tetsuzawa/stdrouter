package stdrouter

import "net/http"

type Router struct{}

func NewRouter() http.Handler { return Router{} }

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request)                         {}
func (router Router) HandleFunc(path interface{}, method interface{}, handlerFunc interface{}) {}
func (router Router) HandleNotFound(handlerFunc interface{})                                   {}
func (router Router) HandleMethodNotAllowed(handlerFunc interface{})                           {}
