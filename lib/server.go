package lib

import (
	"net/http"
)

type ViewHandler func(http.ResponseWriter, *http.Request)

type StrictHandler struct {
	Views           []StrictView
	NotFoundHandler ViewHandler
}

type StrictView struct {
	Patterns []string
	Handler  ViewHandler
}

func NewStrictHandler() *StrictHandler {
	handler := StrictHandler{
		Views: make([]StrictView, 50),
		NotFoundHandler: notFound,
	}
	return &handler
}

func (sh *StrictHandler) HandlePatterns(patterns []string, handler ViewHandler) {
	removeTrailingSlashes(patterns)
	var view StrictView = StrictView{
		Patterns: patterns,
		Handler:  handler,
	}
	sh.Views = append(sh.Views, view)
}

func (sh *StrictHandler) HandlePattern(pattern string, handler ViewHandler) {
	sh.HandlePatterns([]string{pattern}, handler)
}

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

func removeTrailingSlashes(patterns []string) {
	for i, pattern := range patterns {
		patterns[i] = removeTrailingSlash(pattern)
	}
}

func removeTrailingSlash(pattern string) string {
	if string(pattern[len(pattern)-1]) == "/" {
		return pattern[0:len(pattern)-1]
	}
	return pattern
}
