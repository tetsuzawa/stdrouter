package handler

import (
	"fmt"
	"net/http"
)

func GetPosts(w http.ResponseWriter, r *http.Request, userId string) {
	/*
		some implementation ...
	*/
	w.Write([]byte(fmt.Sprintf("get posts. user id: %v", userId)))
}

func GetPost(w http.ResponseWriter, r *http.Request, userId, postId string) {
	/*
		some implementation ...
	*/
	w.Write([]byte(fmt.Sprintf("get post. user id: %v, post id: %v", userId, postId)))
}

func CreatePost(w http.ResponseWriter, r *http.Request, userId string) {
	/*
		some implementation ...
	*/
	w.Write([]byte(fmt.Sprintf("create post. user id: %v", userId)))
}
