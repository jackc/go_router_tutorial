package main

import (
	"fmt"
	"github.com/JackC/go_router_tutorial"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func greetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.FormValue("name"))
}

func main() {
	router := go_router_tutorial.NewRouter()

	router.Get("/", http.HandlerFunc(indexHandler))
	router.Get("/greet/:name", http.HandlerFunc(greetHandler))
	http.Handle("/", router)

	if err := http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
		os.Stderr.WriteString("Could not start web server!\n")
		os.Exit(1)
	}
}
