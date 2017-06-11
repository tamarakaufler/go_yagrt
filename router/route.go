package router

import (
	"errors"
	"fmt"
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

	fmt.Printf(">>> path in getSegments is %+v\n", path)

	// sanity check, that the path starts with a slash
	if ok, err = regexp.MatchString("^/", path); err != nil {
		return nil, err
	}

	if !ok {
		err := errors.New("The path must start with /")
		return nil, err
	}

	segments := strings.Split(path, "/")
	//segments[0] = "/"

	// if the path ends witha slash, we remove the last empty element
	len := len(segments)

	if segments[len-1] == "" {
		segments[len-1] = "/"
	}

	fmt.Printf(">>> segments in getSegments is %+v\n\n", segments[1:len])

	return segments[1:len], nil
}

// processSegment - happens during handler registration
func processSegment(i int, requestPath RequestPath) (map[string]*Route, error) {

	var (
		err error
		ok  bool
	)

	routes := make(map[string]*Route)

	segments := requestPath.Segments
	method := requestPath.Method
	handler := requestPath.Handler

	if len(segments) == 0 {
		return routes, nil
	}

	fmt.Printf(">>> PROCESSING stage %d in segments %v\n", i, segments)
	fmt.Printf("\t>>> processSegment: %+v\n", requestPath.Method)

	seg := segments[i]
	route, ok := routes[seg]

	fmt.Printf("\t>>> processSegment: %+v\n", requestPath.Routes)

	// a new segment => create a new child tree node/Route
	if !ok {
		route = &Route{
			Segment:  seg,
			IsParam:  make(map[string]bool),
			Handlers: make(map[string]Handler),
			Routes:   make(map[string]*Route),
		}
	}

	if ok, err = regexp.MatchString("^:", seg); err != nil {
		return nil, err
	}
	if ok {
		route.IsParam[method] = true
	}

	sl := len(segments)
	fmt.Printf("\t>>> processSegment: Segment = %+v\n", route.Segment)

	//base case of the recursive method
	//associate the path with a handler
	if i == (sl - 1) {
		route.Handlers[method] = handler

		fmt.Printf("\t\t>>> Registering the handler for i=%d : %v\n", i, seg)

		requestPath.Routes[seg] = route

		fmt.Printf("\t>>> processSegment: len(route.Handlers) = %+v\n", len(route.Handlers))
		fmt.Printf("\t>>> processSegment: route.Handlers = %+v\n", route.Handlers)
		fmt.Printf("\t>>> FINISHED (i=%d, segment=%+v) - found the handler\n\n", i, seg)
		return requestPath.Routes, nil
	}

	i++
	seg = requestPath.Segments[i]
	route.Routes[seg] = route
	requestPath.Routes[seg] = route

	fmt.Printf("\t>>> processSegment: recursion for i=%d: segment %s\n", i, requestPath.Segments[i])

	return processSegment(i, requestPath)
}

// processRequest - heppens during HTTP request
//		finds the relevant handler
func processRequest(segments []string, method string, route *Route, i int, params map[string]interface{}) (Handler, map[string]interface{}, error) {

	fmt.Printf("!!! processRequest - parent segment: %+v\nsegments: %+v\ni: %+v\n!!!!!!!!!!\n\n", route.Segment, segments, i)

	var ok bool
	var isParam bool
	var handler Handler

	// base case 1 of the recursive method
	//		path ... /
	if len(segments) == 0 {
		if handler, ok = route.Handlers[method]; !ok {
			return nil, nil, errors.New("No handler for /")
		}
		return handler, params, nil
	}

	// current request path segment
	seg := segments[i]

	// base case 2 of the recursive method
	//		returning when this is the last segment
	if i == (len(segments) - 1) {
		if len(segments) > 1 {
			if handler, ok = route.Handlers[method]; !ok {
				return nil, params, errors.New("No handler for " + fmt.Sprintf("%+v", segments))
			}

			if isParam, ok = route.IsParam[method]; ok && isParam == true {
				params[route.Segment] = seg
			}

			fmt.Printf("\t>>> FINISHED (i=%d, segment=%+v) - found the handler\n\n", i, seg)
			return handler, params, nil
		}
		if route, ok = route.Routes[seg]; !ok {
			return nil, params, errors.New("No route for " + fmt.Sprintf("%v in %+v", seg, segments))

		}
		if handler, ok = route.Handlers[method]; !ok {
			return nil, params, errors.New("No handler for " + fmt.Sprintf("%+v", segments))
		}

		if isParam, ok = route.IsParam[method]; ok && isParam == true {
			params[route.Segment] = seg
		}

		fmt.Printf("\t>>> FINISHED (i=%d, segment=%+v) - found the handler\n\n", i, seg)
		return handler, params, nil
	}

	routes := route.Routes

	if routes == nil {
		return nil, nil, errors.New("Error in route registration: Missing child Routes at " + route.Segment)
	}

	if isParam, ok = route.IsParam[method]; ok && isParam == true {
		params[route.Segment] = seg
	}

	routeSeg := route.Segment

	y := i + 1
	nextSeg := segments[y]

	fmt.Printf("+++ routeSeg=%s, nextSeg=%s\n", routeSeg, nextSeg)
	fmt.Printf("+++ routes=%+v\n\n", routes)

	for key, value := range routes {
		fmt.Printf("###### %v = %+v\n", key, value)
	}

	// recursion into child routes
	if route, ok = routes[nextSeg]; !ok {
		var err error

		// this if seg is a named param
		route, err = getParamSegRoute(routes)

		if err != nil {
			return nil, nil, err
		}
		if route == nil {
			return nil, nil, errors.New("Error in route registration: Missing Route at " + routeSeg)
		}
	}

	return processRequest(segments, method, route, i+1, params)
}

// getParamSegRoute ... returns a Route for a named parameter
func getParamSegRoute(routes map[string]*Route) (*Route, error) {
	var ok bool
	var err error

	for k := range routes {
		ok, err = regexp.MatchString("^:", k)
		if err != nil {
			return nil, err
		}
		if ok {
			return routes[k], nil
		}
	}
	return nil, nil
}
