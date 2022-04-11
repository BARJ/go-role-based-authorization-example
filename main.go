package main

import (
	"log"
	"net/http"
)

func main() {
	log.Printf("Started HTTP server on http://localhost:8888")
	log.Fatalln(http.ListenAndServe(":8888", newRouter()))
}
