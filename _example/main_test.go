package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_newRouter(t *testing.T) {
	// Setup
	r := NewRouter()
	s := httptest.NewServer(r)
	defer s.Close()

	type args struct {
		method  string
		path    string
		reqBody *bytes.Buffer
	}
	type want struct {
		statusCode int
		respBody   string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "/ [get]",
			args: args{
				method: http.MethodGet,
				path:   "/",
			},
			want: want{
				statusCode: http.StatusOK,
				respBody:   "get root",
			},
		},
		{
			name: "/api [get]",
			args: args{
				method: http.MethodGet,
				path:   "/api",
			},
			want: want{
				statusCode: http.StatusOK,
				respBody:   "get api root",
			},
		},
		{
			name: "/api/user",
			args: args{
				method: http.MethodGet,
				path:   "/api/users",
			},
			want: want{
				statusCode: http.StatusOK,
				respBody:   "get users",
			},
		},
		{
			name: "/api/products [get]",
			args: args{
				method: http.MethodGet,
				path:   "/api/products",
			},
			want: want{
				statusCode: http.StatusOK,
				respBody:   "get products",
			},
		},
		{
			name: "/api/users/1 [get]",
			args: args{
				method: http.MethodGet,
				path:   "/api/users/1",
			},
			want: want{
				statusCode: http.StatusOK,
				respBody:   "get user. user id: 1",
			},
		},
		{
			name: "/api/users/1/posts [get]",
			args: args{
				method: http.MethodGet,
				path:   "/api/users/1/posts",
			},
			want: want{
				statusCode: http.StatusOK,
				respBody:   "get posts. user id: 1",
			},
		},
		{
			name: "/api/users/1/profile [get]",
			args: args{
				method: http.MethodGet,
				path:   "/api/users/1/profile",
			},
			want: want{
				statusCode: http.StatusOK,
				respBody:   "get user. user id: 1",
			},
		},
		{
			name: "/api/users/1/posts/5 [get]",
			args: args{
				method: http.MethodGet,
				path:   "/api/users/1/posts/5",
			},
			want: want{
				statusCode: http.StatusOK,
				respBody:   "get post. user id: 1, post id: 5",
			},
		},
		{
			name: "not found [get]",
			args: args{
				method: http.MethodGet,
				path:   "/notfound",
			},
			want: want{
				statusCode: http.StatusNotFound,
				respBody:   "Not Found\n",
			},
		},
		{
			name: "method not allowed / [get]",
			args: args{
				method:  http.MethodDelete,
				path:    "/",
				reqBody: nil,
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				respBody:   "Method Not Allowed\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Request and Assertions
			client := new(http.Client)
			var req *http.Request
			var err error
			switch tt.args.method {
			case http.MethodGet:
				req, err = http.NewRequest(tt.args.method, s.URL+tt.args.path, nil)
			case http.MethodPost:
				req, err = http.NewRequest(tt.args.method, s.URL+tt.args.path, tt.args.reqBody)
			case http.MethodPatch:
				req, err = http.NewRequest(tt.args.method, s.URL+tt.args.path, tt.args.reqBody)
			case http.MethodDelete:
				req, err = http.NewRequest(tt.args.method, s.URL+tt.args.path, nil)
			default:
				t.Fatalf("method not allowed")
			}
			if err != nil {
				t.Fatalf("http.NewRequest: %v", err)
			}
			//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("client.Do: %v", err)
			}

			if !reflect.DeepEqual(resp.StatusCode, tt.want.statusCode) {
				t.Fatalf("resp.StatusCode: got: %d, want: %d", resp.StatusCode, tt.want.statusCode)
			}
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				t.Fatalf("ioutil.ReadAll failed: %s", err)
			}
			got := string(body)
			if !reflect.DeepEqual(got, tt.want.respBody) {
				t.Errorf("request = /%v, got %v, want %v\n", tt.args.path, got, tt.want.respBody)
			}
		})
	}
}
