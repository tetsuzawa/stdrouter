package handler

import "net/http"

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	/*
		some implementation ...
	*/
	http.Error(w, "Not Found", http.StatusNotFound)
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	/*
		some implementation ...
	*/
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
