package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestPostHashHandler(t *testing.T) {

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		code int
	}{
		{"base case", args{httptest.NewRecorder(), httptest.NewRequest("POST", "/hash", strings.NewReader("password=test"))}, http.StatusOK},
		{"base case", args{httptest.NewRecorder(), httptest.NewRequest("POST", "/hash", strings.NewReader("nopassword=test"))}, http.StatusBadRequest},
		{"bad method", args{httptest.NewRecorder(), httptest.NewRequest("GET", "/hash", nil)}, http.StatusMethodNotAllowed},
	}

	var storage = NewStorage()
	var c = NewController(storage)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
			c.PostHash(tt.args.w, tt.args.r)
			w := tt.args.w.(*httptest.ResponseRecorder)
			if got, want := w.Code, tt.code; got != want {
				t.Errorf("got %d; want %d", got, want)
			}
		})
	}
}

func TestGetHashHandler(t *testing.T) {

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		code int
	}{
		{"base case", args{httptest.NewRecorder(), httptest.NewRequest("GET", "/hash/77", nil)}, http.StatusOK},
		{"base case", args{httptest.NewRecorder(), httptest.NewRequest("GET", "/hash/77g", nil)}, http.StatusBadRequest},
		{"base case", args{httptest.NewRecorder(), httptest.NewRequest("GET", "/hash/78", nil)}, http.StatusNotFound},
		{"bad method", args{httptest.NewRecorder(), httptest.NewRequest("POST", "/hash/77", nil)}, http.StatusMethodNotAllowed},
	}

	var storage = NewStorage()
	var c = NewController(storage)
	storage.Set(77, "test_hash")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.GetHash(tt.args.w, tt.args.r)
			w := tt.args.w.(*httptest.ResponseRecorder)
			if got, want := w.Code, tt.code; got != want {
				t.Errorf("got %d; want %d", got, want)
			}
		})
	}
}

func TestHealthHandler(t *testing.T) {

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		code int
	}{
		{"base case", args{httptest.NewRecorder(), httptest.NewRequest("GET", "/health", nil)}, http.StatusOK},
	}
	atomic.StoreInt32(&healthy, 1)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Health(tt.args.w, tt.args.r)
			w := tt.args.w.(*httptest.ResponseRecorder)
			if got, want := w.Code, tt.code; got != want {
				t.Errorf("got %d; want %d", got, want)
			}
		})
	}
}

func TestShutdownHandler(t *testing.T) {

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		code int
	}{
		{"base case", args{httptest.NewRecorder(), httptest.NewRequest("GET", "/shutdown", nil)}, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Shutdown(tt.args.w, tt.args.r)
			w := tt.args.w.(*httptest.ResponseRecorder)
			if got, want := w.Code, tt.code; got != want {
				t.Errorf("got %d; want %d", got, want)
			}
		})
	}
}