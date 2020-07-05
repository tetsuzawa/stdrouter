package main

import (
	"fmt"
	"log"
	"net/http"

	mw "github.com/tetsuzawa/stdrouter/_example/middleware"
)

const (
	host = "127.0.0.1"
	port = "8080"
)

func main() {
	fmt.Println("Server Start....")
	r := NewRouter()
	r = mw.RequestLog(r)
	address := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Server is starting at %s ...", address)
	if err := http.ListenAndServe(address, r); err != nil {
		panic(err)
	}
}
