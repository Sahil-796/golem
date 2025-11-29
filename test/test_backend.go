package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

    // Backend 1 (port 3000)
    go func() {
        mux := http.NewServeMux() // separate router
        mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintln(w, "Response from Backend 1")
        })
        log.Println("Backend 1 running on :3000")
        http.ListenAndServe(":3000", mux)
    }()

    // Backend 2 (port 3001)
    go func() {
        mux := http.NewServeMux() // separate router
        mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintln(w, "Response from Backend 2")
        })
        log.Println("Backend 2 running on :3001")
        http.ListenAndServe(":3001", mux)
    }()

    select {}
}
