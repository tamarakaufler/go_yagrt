package router

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_processRequest(t *testing.T) {

	isParam := make(map[string]bool)
	isParamTrue := make(map[string]bool)
	handlers := make(map[string]Handler)
	handlersEmpty := make(map[string]Handler)
	routes := make(map[string]*Route)
	routesP := make(map[string]*Route)
	routesEmpty := make(map[string]*Route)
	paramsEmpty := make(map[string]interface{})
	params := make(map[string]interface{})

	isParam["GET"] = false
	isParamTrue["GET"] = true

	handler := func(res http.ResponseWriter, req *http.Request) {}
	handlers["GET"] = handler

	routeB := &Route{
		Segment:  "bbb",
		IsParam:  isParam,
		Handlers: handlers,
		Routes:   routesEmpty,
	}
	routes["bbb"] = routeB

	routeP := &Route{
		Segment:  ":bbb",
		IsParam:  isParamTrue,
		Handlers: handlers,
		Routes:   routesEmpty,
	}
	routesP[":bbb"] = routeP

	params[":bbb"] = "101"

	type args struct {
		segments []string
		method   string
		route    *Route
		i        int
		params   map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    Handler
		want1   map[string]interface{}
		wantErr bool
	}{
		{
			name: "/ request - success",
			args: args{
				segments: []string{},
				method:   "GET",
				route: &Route{
					Segment:  "/",
					IsParam:  isParam,
					Handlers: handlers,
					Routes:   routes,
				},
				i:      0,
				params: paramsEmpty,
			},
			want:    handler,
			want1:   paramsEmpty,
			wantErr: false,
		},
		{
			name: "/test request - success",
			args: args{
				segments: []string{"test"},
				method:   "GET",
				route: &Route{
					Segment:  "test",
					IsParam:  isParam,
					Handlers: handlers,
					Routes:   routes,
				},
				i:      0,
				params: paramsEmpty,
			},
			want:    handler,
			want1:   paramsEmpty,
			wantErr: false,
		},
		{
			name: "/aaa/bbb request - success",
			args: args{
				segments: []string{"aaa", "bbb"},
				method:   "GET",
				route: &Route{
					Segment:  "aaa",
					IsParam:  isParam,
					Handlers: handlersEmpty,
					Routes:   routes,
				},
				i:      0,
				params: paramsEmpty,
			},
			want:    routes["bbb"].Handlers["GET"],
			want1:   paramsEmpty,
			wantErr: false,
		},
		{
			name: "/aaa/:bbb request - success",
			args: args{
				segments: []string{"aaa", "101"},
				method:   "GET",
				route: &Route{
					Segment:  "aaa",
					IsParam:  isParam,
					Handlers: handlersEmpty,
					Routes:   routesP,
				},
				i:      0,
				params: paramsEmpty,
			},
			want:    routesP[":bbb"].Handlers["GET"],
			want1:   params,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := processRequest(tt.args.segments, tt.args.method, tt.args.route, tt.args.i, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("processRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.ValueOf(got) != reflect.ValueOf(tt.want) {
				t.Errorf("processRequest() got = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("processRequest() got1 = %+v, want %+v", got1, tt.want1)
			}
		})
	}
}

func Test_processSegment(t *testing.T) {
	isParam := make(map[string]bool)
	isParamTrue := make(map[string]bool)
	isParamTrue["POST"] = true

	handler := func(res http.ResponseWriter, req *http.Request) {}
	handlersA := make(map[string]Handler)
	handlersAB := make(map[string]Handler)
	routes := make(map[string]*Route)
	routesA := make(map[string]*Route)
	routesAB := make(map[string]*Route)
	routesB := make(map[string]*Route)
	routesBP := make(map[string]*Route)

	requestPath1 := RequestPath{
		Segments: []string{"/"},
		Method:   "POST",
		Handler:  handler,
		Routes:   make(map[string]*Route),
	}
	requestPathA := RequestPath{
		Segments: []string{"/", "aaa"},
		Method:   "POST",
		Handler:  handler,
		Routes:   make(map[string]*Route),
	}
	requestPathAB := RequestPath{
		Segments: []string{"/", "aaa", "bbb"},
		Method:   "POST",
		Handler:  handler,
		Routes:   routes,
	}
	requestPathBP := RequestPath{
		Segments: []string{"/", "aaa", ":bbb"},
		Method:   "POST",
		Handler:  handler,
		Routes:   routes,
	}

	handlersA["POST"] = handler
	handlersAB["POST"] = handler

	routeA := &Route{
		Segment:  "aaa",
		IsParam:  isParam,
		Handlers: handlersA,
		Routes:   make(map[string]*Route),
	}
	routesA["aaa"] = routeA

	routeAB := &Route{
		Segment:  "aaa",
		IsParam:  make(map[string]bool),
		Handlers: make(map[string]Handler),
		Routes:   routesAB,
	}
	routeB := &Route{
		Segment:  "bbb",
		IsParam:  isParam,
		Handlers: handlersAB,
		Routes:   routesB,
	}
	routeAB.Routes["bbb"] = routeB
	routesAB["aaa"] = routeAB

	routeBP := &Route{
		Segment:  ":bbb",
		IsParam:  isParamTrue,
		Handlers: handlersAB,
		Routes:   make(map[string]*Route),
	}
	routeA.Routes[":bbb"] = routeBP
	routesBP["aaa"] = routeA

	type args struct {
		i           int
		requestPath RequestPath
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*Route
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "/ url - success",
			args: args{
				i:           0,
				requestPath: requestPath1,
			},
			want:    routes,
			wantErr: false,
		},
		{
			name: "/aaa url - success",
			args: args{
				i:           0,
				requestPath: requestPathA,
			},
			want:    routesA,
			wantErr: false,
		},
		{
			name: "/aaa/bbb url - success",
			args: args{
				i:           0,
				requestPath: requestPathAB,
			},
			want:    routesAB,
			wantErr: false,
		},
		{
			name: "/aaa/:bbb url - success",
			args: args{
				i:           0,
				requestPath: requestPathBP,
			},
			want:    routesBP,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processSegment(tt.args.i, tt.args.requestPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("processSegment() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}

			if a, ok := got[":bbb"]; ok {
				if !reflect.DeepEqual(a.IsParam["POST"], tt.want["aaa"].Routes[":bbb"].IsParam["POST"]) {
					t.Errorf("FAIL processSegment() = %+v, want %+v", got[":bbb"].IsParam["POST"], tt.want["aaa"].Routes[":bbb"].IsParam["POST"])
				}
				if !reflect.DeepEqual(a.Segment, tt.want["aaa"].Routes[":bbb"].Segment) {
					t.Errorf("FAIL processSegment() = %+v, want %+v", got[":bbb"].Segment, tt.want["aaa"].Routes[":bbb"].Segment)
				}

				if reflect.ValueOf(a.Handlers["POST"]) != reflect.ValueOf(tt.want["aaa"].Routes[":bbb"].Handlers["POST"]) {

					t.Errorf("FAIL processSegment() = %+v, want %+v", a.Handlers["POST"], tt.want["aaa"].Routes[":bbb"].Handlers["POST"])
				}

				if reflect.ValueOf(a.Routes["ccc"]) != reflect.ValueOf(tt.want["aaa"].Routes[":bbb"].Routes["ccc"]) {
					t.Errorf("FAIL processSegment() = %+v, want %+v", a.Routes, tt.want["aaa"].Routes[":bbb"].Routes)
				}
			}
		})
	}
}
