package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", helloHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello, World!")
}