package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shellfu/muxer"
	"github.com/shellfu/muxer/middleware"
)

func main() {
	router := muxer.NewRouter()

	router.Use(logRequest, logRequest, middleware.Gzip)
	router.HandlerFunc(http.MethodGet, "/product/:id", func(w http.ResponseWriter, r *http.Request) {
		id := router.Params(r)["id"]
		if _, err := w.Write([]byte("Product ID " + id)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}
	})

	router.HandleRoute(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("index!")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}
	})

	router.HandleRoute(http.MethodGet, "/hello/:name/:last", func(w http.ResponseWriter, r *http.Request) {
		name := router.Params(r)["name"]
		fmt.Println(router.Params(r))
		if _, err := w.Write([]byte("Hello " + name)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}
	})

	router.HandleRoute(http.MethodGet, "/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("danger danger danger!")
	})

	api := router.Subrouter("/api")
	api.HandleRoute(http.MethodGet, "/users/:id", func(w http.ResponseWriter, r *http.Request) {
		id := router.Params(r)["id"]
		if _, err := w.Write([]byte("users " + id)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}
	})

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
