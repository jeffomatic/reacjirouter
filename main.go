package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.ListenAndServe(":1234", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	}))
}
