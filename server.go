package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vaiktorg/grimoire/helpers"
)

func main() {
	router := mux.NewRouter()

	err := Routes(router)
	if err != nil {
		panic(err)
	}

	helpers.ServerCloser(&http.Server{Handler: router, Addr: ":8080"})
}
