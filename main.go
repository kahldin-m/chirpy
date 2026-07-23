package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	serverNew := http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	serverNew.ListenAndServe()

}
