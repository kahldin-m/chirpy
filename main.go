package main

import (
	"net/http"
)

func main() {
	servMux := http.NewServeMux()
	serverNew := http.Server{
		Addr: ":8080",
		Handler: servMux,
	}
	serverNew.ListenAndServe()

}
