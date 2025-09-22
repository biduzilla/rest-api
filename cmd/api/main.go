package api

import (
	"io"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	s.ListenAndServe()
}

func hello(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!")
}
