package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Hi Gustavo, program starting.")
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	server := &http.Server{Addr: ":8080", Handler: mux}
	log.Fatal(server.ListenAndServe())

}
