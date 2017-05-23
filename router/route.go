package router

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
)

// Route - represents configuration for a path segment
type Route struct {
	Segment  string
	IsParam  bool
	Handlers map[string]http.Handler // key=HTTP method, value is the handler
	Routes   map[string]*Route       // represents oprional children segments. key=segment. If no hild Routes => invoke the handler
}

func getSegments(path string) ([]string, error) {
	var ok bool
	var err error

	// sanity chec, thatthe path starts with a slash
	if ok, err = regexp.MatchString("^/", path); err == nil {
		return nil, err
	}
	if !ok {
		err := errors.New("The path must start with /")
		return nil, err
	}

	segments := strings.Split(path, "/")

	// if the path ends witha slash, we remove the last empty element
	len := len(segments)
	if segments[len-1] == "" {
		len = len - 1
	}

	return segments[1:len], nil
}

func processSegment(seg string, routes map[string]*Route, method string, handler http.Handler) (*Route, error) {
	var err error
	var ok bool

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
		return nil, err
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
}
