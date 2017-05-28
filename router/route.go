package router

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strings"
)

// Route - represents configuration for a path segment
type Route struct {
	Segment  string
	IsParam  map[string]bool    // key = method, indicates whether this is a named parameter
	Handlers map[string]Handler // key = HTTP method, value is the handler
	Routes   map[string]*Route  // represents oprional children segments. key = segment. If no hild Routes => invoke the handler
}

func getSegments(path string) ([]string, error) {
	var ok bool
	var err error

	// sanity check, that the path starts with a slash
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

func processSegment(i int, requestPath RequestPath) error {

	var (
		err      error
		ok       bool
		isParam  map[string]bool
		handlers map[string]Handler
		routes   map[string]*Route
	)

	segments := requestPath.Segments
	method := requestPath.Method
	handler := requestPath.Handler

	seg := segments[i]
	route, ok := requestPath.Routes[seg]

	// a new segment => create a new child tree node/Route
	if !ok {
		isParam = make(map[string]bool)
		handlers = make(map[string]Handler)
		routes = make(map[string]*Route)
		route = &Route{
			Segment:  seg,
			IsParam:  isParam,
			Handlers: handlers,
			Routes:   routes,
		}
	}

	if ok, err = regexp.MatchString("^:", seg); err != nil {
		return err
	}
	if ok {
		route.IsParam[method] = true
	}

	sl := len(segments)

	//base case of the recursive method
	//associate the path with a handler
	if i == (sl - 1) {
		route.Handlers[method] = handler
		return nil
	}

	//recursion
	requestPath.Routes[seg] = route

	err = processSegment(i+1, requestPath)
	if err != nil {
		return err
	}

	return nil
}

func processRequest(segments []string, req *http.Request, route *Route, i int) (Handler, error) {

	var ok bool
	var handler Handler

	method := req.Method

	// base case 1 of the recursive method
	//		path ... /
	if len(segments) == 0 {
		if handler, ok = route.Handlers[method]; !ok {
			return nil, errors.New("No handler for " + req.URL.RawPath)
		}
	}

	segment := segments[i]

	if _, ok = route.IsParam[method]; ok {
		ctx := req.Context()
		ctx = context.WithValue(ctx, route.Segment, segment)
		req = req.WithContext(ctx)
	}

	// base case 2 of the recursive method
	//		returning when this is the last segment
	if i == (len(segments) - 1) {
		if handler, ok = route.Handlers[method]; !ok {
			return nil, errors.New("No handler for " + req.URL.RawPath)
		}
		return handler, nil
	}

	routes := route.Routes

	if routes == nil {
		return nil, errors.New("Error in route registration: Missing child Routes at " + route.Segment)
	}

	if route, ok = routes[segment]; !ok {
		return nil, errors.New("Error in route registration: Missing Route at " + route.Segment)
	}

	// recursion into child routes
	processRequest(segments, req, route, i+1)
	return handler, nil
}
