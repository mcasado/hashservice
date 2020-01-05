package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mcasado/hashservice/app"
)

var (
	listenAddr string
)


func main() {
	flag.StringVar(&listenAddr, "listen-addr", "8000", "server listen address")
	flag.Parse()

	// create a logger, router and server
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	router := app.NewRegexRouter()
	server := app.NewServer(
		listenAddr,
		(app.Middlewares{app.Logging(logger), app.Tracing(func() string { return fmt.Sprintf("%d", time.Now().UnixNano()) })}).Apply(router),
		logger,
	)

	// run our server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}