package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/zapkub/cftl/internal/auth"
	"github.com/zapkub/cftl/internal/conf"
	"github.com/zapkub/cftl/internal/database"
	"github.com/zapkub/cftl/internal/frontend"
	"github.com/zapkub/cftl/internal/repository"
	"github.com/zapkub/cftl/internal/sandbox"
)

func main() {

	var mux = http.NewServeMux()

	db, err := database.Open(conf.C.DBDriver, conf.C.DBConnInfo())
	if err != nil {
		log.Fatal(err)
	}
	authenticator := auth.New(repository.New(db))

	var frontendServer = frontend.New(authenticator)
	var sandboxServer sandbox.Server

	frontendServer.Install(mux)
	sandboxServer.Install(mux)

	go func() {
		if err := http.ListenAndServe(conf.C.Address, mux); err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Println("application is now running!!\nvisit http://127.0.0.1:8080")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM|syscall.SIGHUP)
	<-sig
	os.Exit(0)
}
