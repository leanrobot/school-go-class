package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var BASE_TEMPLATE = "templates/base.html"

// simple type to alias the viewHandler interface in net/http.
type ViewHandler func(http.ResponseWriter, *http.Request)

// Strict Handler conforms to the Handler interface in net/http.
// https://golang.org/pkg/net/http/#Handler
type StrictHandler struct {
	Views           []StrictView
	NotFoundHandler ViewHandler
}

// StrictView is a simple association between url resource paths and
// a handler view.
type StrictView struct {
	Patterns []string
	Handler  ViewHandler
}

/*
NewStrictHandler is the constructor for the StrictHandler type.
No custom 404 handler is declared here. to set it change the NotFoundHandler
value.

	strictView := lib.NewStrictHandler()
	strictView.NotFoundHandler = <type func(http.ResponseWriter, *http.Request)>
*/
func NewStrictHandler() *StrictHandler {
	handler := StrictHandler{
		Views: make([]StrictView, 50),
	}
	return &handler
}

/*
HandlePatterns registers several url resources to a view. Trailing slashes are
sanitized out of the url patterns automatically.
*/
func (sh *StrictHandler) HandlePatterns(patterns []string, handler ViewHandler) {
	removeTrailingSlashes(patterns)
	var view StrictView = StrictView{
		Patterns: patterns,
		Handler:  handler,
	}
	sh.Views = append(sh.Views, view)
	fmt.Fprintf(os.Stderr, "%v registered\n", patterns)

}

/*
HandlePattern registers a single url resource to a view. Please reference the
documentation for HandlePatterns for more information.
*/
func (sh *StrictHandler) HandlePattern(pattern string, handler ViewHandler) {
	sh.HandlePatterns([]string{pattern}, handler)
}

func (sh *StrictHandler) ServeStaticFile(pattern string, filename string) {
	sh.HandlePattern(pattern,
		func(res http.ResponseWriter, req *http.Request) {
			buf, err := ioutil.ReadFile(filename)
			if err != nil {
				sh.NotFoundHandler(res, req)
			} else {
				res.Write(buf)
			}
		})
}

/*
ServeHTTP conforms to the http.Handler interface type.

Given a request it determines which view to call for that url resource.
Url patterns are stored in a list to maintain the register-first match-first
relationship between resource and views.
*/
func (sh *StrictHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	url := removeTrailingSlash(req.URL.Path)
	for _, view := range sh.Views {
		for _, pattern := range view.Patterns {
			if pattern == url {
				view.Handler(res, req)
				return
			}
		}
	}

	// Handle 404
	if sh.NotFoundHandler != nil {
		sh.NotFoundHandler(res, req)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}

//TODO: write tests for trailing / removal

// removeTrailingSlashes calls removeTrailingSlash for all strings in the slice.
func removeTrailingSlashes(patterns []string) {
	for i, pattern := range patterns {
		patterns[i] = removeTrailingSlash(pattern)
	}
}

// removeTrailingSlash will return the passed pattern, truncating off a trailing
// slash if it appears in the pattern.
func removeTrailingSlash(pattern string) string {
	if string(pattern[len(pattern)-1]) == "/" {
		fmt.Fprintln(os.Stderr, "trailingSlash removed from", pattern)
		return pattern[0 : len(pattern)-1]
	}
	return pattern
}
