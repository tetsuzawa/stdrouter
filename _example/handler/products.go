package handler

import "net/http"

func GetProducts(w http.ResponseWriter, r *http.Request) {
	/*
		some implementation ...
	*/
	w.Write([]byte("get products"))
}

func CreateProducts(w http.ResponseWriter, r *http.Request) {
	/*
		some implementation ...
	*/
	w.Write([]byte("create products"))
}
