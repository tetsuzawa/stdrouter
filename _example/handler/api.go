package handler

import "net/http"

func GetAPIRoot(w http.ResponseWriter, r *http.Request) {
	/*
		some implementation ...
	*/
	w.Write([]byte("get api root"))
}
