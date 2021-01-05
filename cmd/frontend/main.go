package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/zapkub/cftl/internal/frontend"
	"github.com/zapkub/cftl/internal/sandbox"
)

func main() {

	var mux = http.NewServeMux()

	var frontendServer frontend.Server
	var sandboxServer sandbox.Server

	frontendServer.Install(mux)
	sandboxServer.Install(mux)

	go func() {
		if err := http.ListenAndServe("127.0.0.1:8080", mux); err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("application is now running!! visit http://127.0.0.1:8080")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM|syscall.SIGHUP)
	<-sig
	os.Exit(0)
}
