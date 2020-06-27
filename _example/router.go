// Copyright (c) 2020 Tetsu Takizawa

//+build stdrouter

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"net/http"

	"github.com/tetsuzawa/stdrouter"
	"github.com/tetsuzawa/stdrouter/_example/handler"
)

// NewHandler creates a http router. It passes HTTP requests to the function.
func NewRouter() http.Handler {
	r := stdrouter.NewRouter()
	r.HandleFunc("/", http.MethodGet, handler.RootGET)
	r.HandleFunc("/api", http.MethodGet, handler.API)
	r.HandleFunc("/api/users", http.MethodGet, handler.GetUsers)
	r.HandleFunc("/api/products", http.MethodGet, handler.GetProducts)
	r.HandleFunc("/api/products", http.MethodPost, handler.CreateProducts)
	r.HandleFunc("/api/users/create", http.MethodPost, handler.CreateUser)
	r.HandleFunc("/api/users/:user_id", http.MethodGet, handler.GetUser)
	r.HandleFunc("/api/users/:user_id", http.MethodPatch, handler.UpdateUser)
	r.HandleFunc("/api/users/:user_id", http.MethodDelete, handler.DeleteUser)
	r.HandleFunc("/api/users/:user_id/posts", http.MethodGet, handler.GetPosts)
	r.HandleFunc("/api/users/:user_id/profile", http.MethodGet, handler.GetUser)
	r.HandleFunc("/api/users/:user_id/posts/:post_id", http.MethodGet, handler.GetPost)
	r.HandleFunc("/api/users/:user_id/posts/:post_id/aaa", http.MethodGet, handler.GetPost)
	r.HandleFunc("/api/users/:user_id/posts/:post_id/aaa/bbb", http.MethodGet, handler.GetPost)
	r.HandleNotFound(handler.MethodNotAllowedHandler)
	r.HandleMethodNotAllowed(handler.MethodNotAllowedHandler)
	/*
		...
	*/
	return r
}
