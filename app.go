package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/mux"
	"github.com/tj/go-config"
)

type Options struct {
	Port int `help:"port number"`
}

type String string

func (s String) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, s)
}

func main() {
	options := Options{
		Port: 8080,
	}
	config.MustResolve(&options)

	port := options.Port
	s := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: newHandler(),
	}
	gracehttp.Serve(s)
}

func newHandler() http.Handler {
	r := mux.NewRouter()
	r.Handle("/", String("hello")).Methods("GET")
	r.HandleFunc("/json", myHandler(func(w http.ResponseWriter, r *http.Request) {
		obj := map[string]string{}
		obj["key"] = "val"
		json.NewEncoder(w).Encode(obj)
	})).Methods("GET")
	return r
}

func myHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	// Add filtter
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}
