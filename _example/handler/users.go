package handler

import (
	"fmt"
	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	/*
		some implementation ...
	*/
	w.Write([]byte("get users"))
}

func GetUser(w http.ResponseWriter, r *http.Request, userId string) {
	/*
		some implementation ...
	*/
	w.Write([]byte(fmt.Sprintf("get user. user id: %v", userId)))
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	/*
		some implementation ...
	*/
	w.Write([]byte(fmt.Sprintf("create user")))
}

func UpdateUser(w http.ResponseWriter, r *http.Request, userId string) {
	/*
		some implementation ...
	*/
	w.Write([]byte(fmt.Sprintf("update user. user id: %v", userId)))
}

func DeleteUser(w http.ResponseWriter, r *http.Request, userId string) {
	/*
		some implementation ...
	*/
	w.Write([]byte(fmt.Sprintf("delete user. user id: %v", userId)))
}
