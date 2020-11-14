package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8002",
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/hello", hello)

	log.Println("ListenAndServe")
	server.ListenAndServe()
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world.")
}
