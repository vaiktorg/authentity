package main

import (
	"context"
	"fmt"
	"github.com/vaiktorg/authentity"
	"github.com/vaiktorg/authentity/handlers"
	"os"
	"os/signal"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	err := handlers.Routes(router)
	if err != nil {
		panic(err)
	}

	fmt.Printf("authentity Instance: %s\n", authentity.Global.Issuer)

	ServerCloser(&http.Server{Handler: router, Addr: ":8080"})
}

func ServerCloser(server *http.Server) {
	go func() {
		err := server.ListenAndServe()
		if err == os.ErrClosed {
			fmt.Println(err)
		}
	}()
	fmt.Println("Listening on", server.Addr)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	err := server.Shutdown(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Shutting Down")
	os.Exit(0)
}
