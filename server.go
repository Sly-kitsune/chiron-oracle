package main

import (
	"log"
	"mime"
	"net/http"
	"os"
)

func main() {
	// 1. Force register WASM mime type (essential for Alpine/Cloudflare)
	mime.AddExtensionType(".wasm", "application/wasm")
	mime.AddExtensionType(".js", "application/javascript")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	log.Printf("Server starting on :%s...", port)
	log.Println("Serving files from ./public")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
