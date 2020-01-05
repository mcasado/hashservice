package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestApplyMiddleware(t *testing.T) {
	h := func(http.ResponseWriter, *http.Request) {}
	type args struct {
		h           http.HandlerFunc
		middlewares Middlewares
	}

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	tests := []struct {
		name string
		args args
	}{
		{"base-case", args{h, []Middleware{Logging(logger)}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.middlewares.Apply(tt.args.h)
		})
	}
}

func TestLogger(t *testing.T) {
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("testing")
	})
	//h := func(http.ResponseWriter, *http.Request) {}
	type args struct {
		l *log.Logger
	}
	tests := []struct {
		name string
		args args
	}{
		{"base-case", args{log.New(os.Stdout, "http: ", log.LstdFlags)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Logging(tt.args.l)(f)
		})
	}
}
