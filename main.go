package main

import (
	"fmt"
	"net/http"

	"github.com/tamarakaufler/go_router/router"
)

func getHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(200)
	res.Write([]byte("Hello World"))
}

func main() {
	fmt.Println("go_router wheel")

	mux, err := router.New("/")
	if err != nil {
		panic(err)
	}

	mux.GET("/", getHandler)
	mux.GET("/test", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte("Hello Test without slash"))
	})
	mux.GET("/test/", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte("Hello Test with slash"))
	})
	mux.GET("/aaa/bbb", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte("Hello aaa/bbb"))
	})
	mux.GET("/aaa/bbb/ccc", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte("Hello aaa/bbb/ccc"))
	})
	mux.GET("/param/:param", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)

		param := req.Context().Value(":param")
		fmt.Printf("%v\n\n", param)
		res.Write([]byte("Hello Param"))
	})

	http.ListenAndServe(":8888", mux)
}
