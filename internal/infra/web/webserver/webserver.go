package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HandlerMethod struct {
	Route   string
	Method  string
	Handler http.HandlerFunc
}

func NewHandlerMethod(route string, method string, handler http.HandlerFunc) *HandlerMethod {
	return &HandlerMethod{
		Route:   route,
		Method:  method,
		Handler: handler,
	}
}

type WebServer struct {
	Router        chi.Router
	Handlers      []HandlerMethod
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make([]HandlerMethod, 0), // []HandlerMethod{}
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(handler HandlerMethod) {
	s.Handlers = append(s.Handlers, handler)
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for _, handler := range s.Handlers {
		s.Router.Method(handler.Method, handler.Route, handler.Handler)
	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}
