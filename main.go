package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	listenAddr string
)


func main() {
	flag.StringVar(&listenAddr, "listen-addr", "8000", "server listen address")
	flag.Parse()

	// create a logger, router and server
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	router := newRegexRouter()
	server := newServer(
		listenAddr,
		(Middlewares{Logging(logger), Tracing(func() string { return fmt.Sprintf("%d", time.Now().UnixNano()) })}).Apply(router),
		logger,
	)

	// run our server
	if err := server.run(); err != nil {
		log.Fatal(err)
	}
}