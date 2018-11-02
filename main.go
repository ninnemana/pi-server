package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ninnemana/pi-server/handlers"
)

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Write([]byte("This is an example server.\n"))
	})

	mux.HandleFunc("/webhooks/github", handlers.Github)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("fell out of serving traffic: %v", err)
	}
}
