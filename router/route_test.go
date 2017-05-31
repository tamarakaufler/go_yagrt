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
		// TODO: Add test cases.
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processSegment(tt.args.i, tt.args.requestPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("processSegment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processSegment() = %v, want %v", got, tt.want)
			}
		})
	}
}
