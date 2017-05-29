package router

import (
	"context"
	"errors"
	"fmt"
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

	fmt.Printf(">>> path in getSegments is %v\n", path)

	// sanity check, that the path starts with a slash
	if ok, err = regexp.MatchString("^/", path); err != nil {
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

	fmt.Printf(">>> segments in getSegments is %v\n\n", segments[1:len])

	return segments[1:len], nil
}

// processSegment - happens during handler registration
func processSegment(i int, requestPath RequestPath) (map[string]*Route, error) {

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
	routes = requestPath.Routes

	if len(segments) == 0 {
		return routes, nil
	}

	fmt.Printf(">>> processSegment: i=%d, segments %v\n", i, segments)
	fmt.Printf("\t>>> processSegment: +%v\n", requestPath.Method)

	seg := segments[i]
	route, ok := routes[seg]

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
		return nil, err
	}
	if ok {
		route.IsParam[method] = true
	}

	sl := len(segments)
	fmt.Printf("\t>>> processSegment: Segment = %v\n", route.Segment)

	//base case of the recursive method
	//associate the path with a handler
	if i == (sl - 1) {
		route.Handlers[method] = handler
		routes[seg] = route

		//fmt.Printf("\t>>> processSegment: handler = %v\n")
		fmt.Printf("\t>>> processSegment: len(route.Handlers) = %v\n", len(route.Handlers))
		fmt.Printf("\t>>> FINISHED (i=%d, segment=%v) - found the handler\n\n", i, seg)
		return routes, nil
	}

	routes[seg] = route
	requestPath.Routes = routes

	fmt.Printf("\t>>> processSegment: recursion for segment %s\n", requestPath.Segments[i+1])

	processSegment(i+1, requestPath)

	return routes, nil
}

// processRequest - heppens during HTTP request
//		finds the relevant handler
func processRequest(segments []string, req *http.Request, route *Route, i int) (Handler, error) {

	fmt.Printf(">>> processRequest - segments: %v\nroute: %v\ni: %v\n", segments, route, i)

	var ok bool
	var handler Handler

	method := req.Method

	fmt.Printf("\t>>> processRequest - Routes length=%d\n", len(segments))

	// base case 1 of the recursive method
	//		path ... /
	if len(segments) == 0 {
		if handler, ok = route.Handlers[method]; !ok {
			return nil, errors.New("No handler for " + req.URL.RawPath)
		}
		return handler, nil
	}

	routes := route.Routes
	segment := segments[i]

	if routes == nil {
		return nil, errors.New("Error in route registration: Missing child Routes at " + route.Segment)
	}

	fmt.Printf("\t>>> processRequest - route for segment %s=%v\n", segment, routes[segment])

	if route, ok = routes[segment]; !ok {
		return nil, errors.New("Error in route registration: Missing Route at " + route.Segment)
	}

	if _, ok = route.IsParam[method]; ok {
		ctx := req.Context()
		ctx = context.WithValue(ctx, route.Segment, segment)
		req = req.WithContext(ctx)

		fmt.Printf("\t\t>> processRequest -  parameter %v = %v\n\n", route.Segment, req.Context().Value(route.Segment))
	}

	// base case 2 of the recursive method
	//		returning when this is the last segment
	if i == (len(segments) - 1) {
		if handler, ok = route.Handlers[method]; !ok {
			return nil, errors.New("No handler for " + req.URL.RawPath)
		}

		fmt.Printf("\t>>> FINISHED (i=%d, segment=%v) - found the handler\n\n", i, segment)
		return handler, nil
	}

	// recursion into child routes
	processRequest(segments, req, route, i+1)
	return handler, nil
}
