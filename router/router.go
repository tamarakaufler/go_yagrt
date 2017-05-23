package router

import (
	"context"
	"net/http"
	"regexp"
)

// Router - Router object needs to satisfy this interface to be a handler
type Router interface {
	http.Handler // ie implements the ServeHTTP method, writes the response out to the client
	GET(path string, handler http.Handler)
	POST(path string, handler http.Handler)
	DELETE(path string, handler http.Handler)
	PUT(path string, handler http.Handler)
}

// Mux object
type Mux struct {
	Context     context.Context
	RootSegment string
	Params      map[string]interface{} //?
	Routes      map[string]*Route
}

// ServeHTTP method uses the paths and their handlers
// to execure the correct handler for the request URL
func (m *Mux) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// split the req.Url.Path
	/*
		segments, err := getSegments(req.URL.Path)
		if err != nil {
			panic(e)
		}
	*/
	// create the Routes (Route trie tree)
	// extract the parameters, if any => attach as Params
	// get the handler and execute it
}

// New method for creating a new router, ie a new Mux instance
// after creating a new Mux, paths need to be registere
// (through GET/POST/PUT/DELETE methods)
// with all the relevant information like:
//			HTTP method
//			handler etc
func New(rootSegment string) *Mux {
	if rootSegment == "" {
		rootSegment = "/"
	}

	handlers := make(map[string]http.Handler)
	routes := make(map[string]*Route)

	baseRoute := &Route{
		Segment:  rootSegment,
		IsParam:  false,
		Handlers: handlers,
		Routes:   routes,
	}

	routes[rootSegment] = baseRoute

	m := &Mux{
		RootSegment: rootSegment,
		Routes:      routes,
	}

	return m
}

// GET method
func (m *Mux) GET(path string, handler http.Handler) {
}

// POST method
func (m *Mux) POST(path string, handler http.Handler) {
}

// PUT method
func (m *Mux) PUT(path string, handler http.Handler) {
}

// DELETE method
func (m *Mux) DELETE(path string, handler http.Handler) {
}

// ListenAndServe method - we arproviding our custom router
func (m *Mux) ListenAndServe(port string) error {
	return http.ListenAndServe(port, m)
}

// Private methods and functions
//----------------------------------------------------------------------------------

func (m *Mux) register(method string, segments []string, handler http.Handler) error {

	// only root path registered
	if len(segments) == 0 {
		m.Routes[m.RootSegment].Handlers[method] = handler

		return nil
	}

	//register the path by inserting it into the Route tree
	var err error

	routes := m.Routes

	for i, seg := range segments {
		route, ok := routes[seg]

		// a new segment => create a new child tree node/Route
		if !ok {
			handlers := make(map[string]http.Handler)
			routes = make(map[string]*Route)
			route = &Route{
				Segment:  seg,
				IsParam:  false,
				Handlers: handlers,
				Routes:   routes,
			}
		}

		if ok, err = regexp.MatchString("^:", seg); err != nil {
			return err
		}
		if ok {
			route.IsParam = true
		}

		routes = route.Routes

		sl := len(segments)

		//associate the path with a handler
		if i == (sl - 1) {
			route.Handlers[method] = handler
		}

		routes[method] = route
		route.Routes = routes

		//segRoutes := make(map[string]*Route)
	}

	return nil
}
