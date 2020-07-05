package handler

import "net/http"

func GetRoot(w http.ResponseWriter, r *http.Request) {
	/*
		some implementation ...
	*/
	w.Write([]byte("get root"))
}
