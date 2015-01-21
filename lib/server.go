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
		// NotFoundHandler: notFound,
	}
	return &handler
}

func (sh *StrictHandler) HandlePatterns(patterns []string, handler ViewHandler) {
	var view StrictView = StrictView{
		Patterns: patterns,
		Handler:  handler,
	}
	sh.Views = append(sh.Views, view)
}

func (sh *StrictHandler) HandlePattern(pattern string, handler ViewHandler) {
	sh.HandlePatterns([]string{pattern}, handler)
}

// Handle(["/", "/index.html"], indexHTML)

func (sh *StrictHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	url := req.URL.Path
	for _, view := range sh.Views {
		for _, pattern := range view.Patterns {
			if pattern == url {
				view.Handler(res, req)
				return
			}
		}
	}
	//TODO(assign2): add handling for having no NotFoundHandler
	//sh.NotFoundHandler(res, req)
}
