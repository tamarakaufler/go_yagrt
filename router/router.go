package router

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/davecgh/go-spew/spew"
)

// Handler type
type Handler func(http.ResponseWriter, *http.Request)

// Router - Router object needs to satisfy this interface to be a handler
type Router interface {
	ServeHTTP(http.ResponseWriter, *http.Request) // ie implements the ServeHTTP method, writes the response out to the client
	GET(path string, handler Handler)
	POST(path string, handler Handler)
	DELETE(path string, handler Handler)
	PUT(path string, handler Handler)
}

// Mux object
type Mux struct {
	BaseRoute *Route
}

// RequestPath type
type RequestPath struct {
	Segments []string
	Routes   map[string]*Route
	Method   string
	Handler  Handler
}

// New method for creating a new router, ie a new Mux instance
// after creating a new Mux, paths need to be registere
// (through GET/POST/PUT/DELETE method)
// with all the relevant information like:
//			HTTP method
//			handler etc
func New(rootSegment string) (*Mux, error) {
	var ok bool
	var err error
	if ok, err = regexp.MatchString("^/", rootSegment); err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("base root must start with /")
	}
	isParam := make(map[string]bool)
	handlers := make(map[string]Handler)
	routes := make(map[string]*Route)

	baseRoute := &Route{
		Segment:  rootSegment,
		IsParam:  isParam,
		Handlers: handlers,
		Routes:   routes,
	}

	m := &Mux{
		BaseRoute: baseRoute,
	}

	return m, nil
}

// ServeHTTP method uses the paths and their handlers
// to execure the correct handler for the request URL
func (m *Mux) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	// extract the parameters, if any => add as a context parameter
	// get the relevant handler
	// users of this package, using named parameters
	// have access to them through the request context and the
	// parameter name used when registering the url, eg:
	// 		req.Context().Value("Username");
	segments, err := getSegments(req.URL.Path)
	if err != nil {
		log.Fatal(err)
	}

	handler, err := processRequest(segments, req, m.BaseRoute, 0)
	if err != nil {
		notFound(res, req)
	}

	handler(res, req)
}

// GET method
func (m *Mux) GET(path string, handler Handler) {
	fmt.Printf("---> %v - %v\n", path, handler)
	m.register(path, "GET", handler)
}

// POST method
func (m *Mux) POST(path string, handler Handler) {
	m.register(path, "POST", handler)
}

// PUT method
func (m *Mux) PUT(path string, handler Handler) {
	m.register(path, "PUT", handler)
}

// DELETE method
func (m *Mux) DELETE(path string, handler Handler) {
	m.register(path, "DELETE", handler)
}

// ListenAndServe method - we are providing our custom router
func (m *Mux) ListenAndServe(port string) error {
	return http.ListenAndServe(port, m)
}

// Private methods and functions
//----------------------------------------------------------------------------------
func notFound(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(404)
	res.Write([]byte("Not found"))
}

func (m *Mux) register(path string, method string, handler Handler) {
	segments, err := getSegments(path)

	if err != nil {
		fmt.Println("=====================================")
		panic(err)
	}

	m.doRegister(method, segments, handler)
}

func (m *Mux) doRegister(method string, segments []string, handler Handler) error {

	// only root path registered
	if len(segments) == 0 {
		m.BaseRoute.Handlers[method] = handler

		return nil
	}

	// all information needed to register a handler
	requestPath := RequestPath{
		Segments: segments,           // decomposed path
		Routes:   m.BaseRoute.Routes, // parent tree child routes
		Method:   method,
		Handler:  handler, // handler to be registered for this url
	}

	routes, err := processSegment(0, requestPath)
	if err != nil {
		return err
	}

	m.BaseRoute.Routes = routes

	if x, ok := m.BaseRoute.Routes["bbb"]; ok {
		spew.Dump(x.Segment)
		spew.Dump(x.IsParam["GET"])
		spew.Dump(x.Handlers["GET"])
		spew.Dump(x.Routes["GET"])
	}

	return nil
}
