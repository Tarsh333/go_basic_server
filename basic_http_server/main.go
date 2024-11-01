package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("url: %s",r.RequestURI)
	fmt.Fprint(w, "  ", r.Host, "  ", r.Method, "  ", r.RequestURI)
}
func main() {
	port := 8000
	http.HandleFunc("/", handler)
	fmt.Printf("server started at port %d", port)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
