package handler

import "net/http"

func RootGET(w http.ResponseWriter, r *http.Request) {
	/*
		some implementation ...
	*/
	w.Write([]byte("get root"))
}

func RootPOST(w http.ResponseWriter, r *http.Request) {
	/*
		some implementation ...
	*/
	w.Write([]byte("post root"))
}
