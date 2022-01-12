package main

import (
	"fmt"
	"github.com/vaiktorg/Authentity"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vaiktorg/grimoire/helpers"
)

func main() {
	router := mux.NewRouter()

	err := authentity.Routes(router)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Authentity Instance: %s\n", authentity.Global.Issuer)

	helpers.ServerCloser(&http.Server{Handler: router, Addr: ":8080"})
}
