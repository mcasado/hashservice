package main

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"syscall"
	"time"

	"io"
	"net/http"
	"sync/atomic"
)

func (c *Controller) PostHash(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	password := r.FormValue("password")
	if password == "" {
		log.Println("Password was not provided")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "Password was not provided")
		return
	}

	identifier := c.storage.IncrementIdentifier()
	go func(id int64) {
		time.Sleep(5 * time.Second)
		hash := CreateHash(password)
		c.storage.Set(id, hash)
		if err := Save("./file.tmp", c.storage.Map()); err != nil {
			log.Fatalln(err)
		}
	}(identifier)

	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, strconv.FormatInt(identifier, 10))
}

func (c *Controller) GetHash(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/hash/")
	if n, err := strconv.Atoi(id); err == nil {
		hash := c.storage.Get(int64(n))
		if hash != "" {
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, hash)
		} else {
			w.WriteHeader(http.StatusNotFound)
			_, _ = io.WriteString(w, id+" has no matching hash")
		}

	} else {
		log.Println(id, "is not an integer.")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, id+" is not an integer.")
		return
	}
}

func Health(w http.ResponseWriter, r *http.Request) {

	if atomic.LoadInt32(&healthy) == 1 {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = io.WriteString(w, `{"alive": true}`)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

func Shutdown(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("shutting down ..."))
	// send SIGTERM to quit channel to start shutdown
	quit <- syscall.SIGTERM
}

func Stats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(statsMiddleware.Data())
	w.Write(b)
}
