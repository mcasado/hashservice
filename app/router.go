package app

import (
	"log"
	"net/http"
	"regexp"
)

var storage = NewStorage()
var c = NewController(storage)

func init() {

	if err := Load("./file.tmp", &storage.hashes); err != nil {
		log.Println(err)
	}

	// To store the keys in slice in sorted order
	keys := make([]int, len(storage.hashes))
	for k := range storage.hashes {
		keys = append(keys, int(k))
	}

	storage.hashIdentifier = int64(Max(keys))
}

func NewRegexRouter() *regexResolver {

	rr := newPathResolver()
	rr.Add("GET /health",  Health)
	rr.Add("GET /hash(/\\d+\\z)", c.GetHash)
	rr.Add("POST /hash\\z", c.PostHash)
	rr.Add("GET /shutdown", Shutdown)
	rr.Add("GET /stats", Stats)
	// routes
	return rr
}

func newPathResolver() *regexResolver {
	return &regexResolver{
		handlers: make(map[string]http.HandlerFunc),
		cache:    make(map[string]*regexp.Regexp),
	}
}

type regexResolver struct {
	handlers map[string]http.HandlerFunc
	cache    map[string]*regexp.Regexp
}

func (r *regexResolver) Add(regex string, handler http.HandlerFunc) {
	r.handlers[regex] = handler
	cache, _ := regexp.Compile(regex)
	r.cache[regex] = cache
}

func (r *regexResolver) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	check := req.Method + " " + req.URL.Path
	for pattern, handlerFunc := range r.handlers {
		if r.cache[pattern].MatchString(check) == true {
			handlerFunc(res, req)
			return
		}
	}

	http.NotFound(res, req)
}