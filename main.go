package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	whisperTarget, err := url.Parse("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	whisperProxy := httputil.NewSingleHostReverseProxy(whisperTarget)

	defaultTarget, err := url.Parse("http://localhost:8880")
	if err != nil {
		log.Fatal(err)
	}
	defaultProxy := httputil.NewSingleHostReverseProxy(defaultTarget)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/audio/transcriptions", "/whisper-metrics":
			if r.URL.Path == "/whisper-metrics" {
				r.URL.Path = "/metrics"
			}
			whisperProxy.ServeHTTP(w, r)
		default:
			defaultProxy.ServeHTTP(w, r)
		}
	})

	log.Println("Starting proxy server")
	if err := http.ListenAndServe(":8089", nil); err != nil {
		log.Fatal(err)
	}
}
