package stdrouter

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"
)

func (n *Node) registerParent() {
	Walk(n, func(node *Node) bool {
		for _, c := range node.Children {
			c.Parent = node
		}
		return true
	})
}

func TestNode_Add(t *testing.T) {
	type fields struct {
		Depth       uint
		Endpoint    string
		IsPathParam bool
		Methods     map[string]HandlerFunc
		Parent      *Node
		Children    []*Node
	}
	type args struct {
		p           string
		httpMethod  string
		handlerFunc HandlerFunc
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		wantField fields
	}{
		{
			name:   "Add root to empty tree",
			fields: fields{},
			args: args{
				p:           "/",
				httpMethod:  http.MethodGet,
				handlerFunc: HandlerFunc{Package: "handler", Func: "GetRoot"},
			},
			wantErr: false,
			wantField: fields{
				Depth:    0,
				Endpoint: "",
				Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
				Children: nil,
			},
		},
		{
			name: "Add api to rooted tree",
			fields: fields{
				Depth:    0,
				Endpoint: "",
				Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
				Children: nil,
			},
			args: args{
				p:           "/api",
				httpMethod:  http.MethodGet,
				handlerFunc: HandlerFunc{Package: "handler", Func: "GetAPI"},
			},
			wantErr: false,
			wantField: fields{
				Depth:    0,
				Endpoint: "",
				Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
				Children: []*Node{
					{
						Depth:    1,
						Endpoint: "api",
						Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetAPI"}},
						Children: nil,
					},
				},
			},
		},
		{
			name:   "Add api to empty tree",
			fields: fields{},
			args: args{
				p:           "/api",
				httpMethod:  http.MethodGet,
				handlerFunc: HandlerFunc{Package: "handler", Func: "GetAPI"},
			},
			wantErr: false,
			wantField: fields{
				Depth:    0,
				Endpoint: "",
				Methods:  nil,
				Children: []*Node{
					{
						Depth:    1,
						Endpoint: "api",
						Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetAPI"}},
					},
				},
			},
		},
		{
			name: "Add [GET] users to tree",
			fields: fields{
				Depth:    0,
				Endpoint: "",
				Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
				Children: []*Node{
					{
						Depth:    1,
						Endpoint: "api",
						Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetAPI"}},
						Children: nil,
					},
				},
			},
			args: args{
				p:           "/api/users",
				httpMethod:  http.MethodGet,
				handlerFunc: HandlerFunc{Package: "handler", Func: "GetUsers"},
			},
			wantErr: false,
			wantField: fields{
				Depth:    0,
				Endpoint: "",
				Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
				Children: []*Node{
					{
						Depth:    1,
						Endpoint: "api",
						Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetAPI"}},
						Children: []*Node{
							{
								Depth:    2,
								Endpoint: "users",
								Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetUsers"}},
								Children: nil,
							},
						},
					},
				},
			},
		},
		{
			name: "Add [POST] users to tree",
			fields: fields{
				Depth:    0,
				Endpoint: "",
				Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
				Children: []*Node{
					{
						Depth:    1,
						Endpoint: "api",
						Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetAPI"}},
						Children: []*Node{
							{
								Depth:    2,
								Endpoint: "users",
								Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetUsers"}},
								Children: nil,
							},
						},
					},
				},
			},
			args: args{
				p:           "/api/users",
				httpMethod:  http.MethodPost,
				handlerFunc: HandlerFunc{Package: "handler", Func: "CreateUsers"},
			},
			wantErr: false,
			wantField: fields{
				Depth:    0,
				Endpoint: "",
				Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
				Children: []*Node{
					{
						Depth:    1,
						Endpoint: "api",
						Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetAPI"}},
						Children: []*Node{
							{
								Depth:    2,
								Endpoint: "users",
								Methods: map[string]HandlerFunc{
									http.MethodGet:  {"handler", "GetUsers"},
									http.MethodPost: {"handler", "CreateUsers"},
								},
								Children: nil,
							},
						},
					},
				},
			},
		},
		{
			name: "Add [GET] products to tree",
			fields: fields{
				Depth:    0,
				Endpoint: "",
				Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
				Children: []*Node{
					{
						Depth:    1,
						Endpoint: "api",
						Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetAPI"}},
						Children: []*Node{
							{
								Depth:    2,
								Endpoint: "users",
								Methods: map[string]HandlerFunc{
									http.MethodGet:  {"handler", "GetUsers"},
									http.MethodPost: {"handler", "CreateUsers"},
								},
								Children: nil,
							},
						},
					},
				},
			},
			args: args{
				p:           "/api/products",
				httpMethod:  http.MethodGet,
				handlerFunc: HandlerFunc{Package: "handler", Func: "GetProducts"},
			},
			wantErr: false,
			wantField: fields{
				Depth:    0,
				Endpoint: "",
				Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
				Children: []*Node{
					{
						Depth:    1,
						Endpoint: "api",
						Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetAPI"}},
						Children: []*Node{
							{
								Depth:    2,
								Endpoint: "users",
								Methods: map[string]HandlerFunc{
									http.MethodGet:  {"handler", "GetUsers"},
									http.MethodPost: {"handler", "CreateUsers"},
								},
								Children: nil,
							},
							{
								Depth:    2,
								Endpoint: "products",
								Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetProducts"}},
								Children: nil,
							},
						},
					},
				},
			},
		},
		{
			name: "Add create users to tree",
			fields: fields{
				Depth:    0,
				Endpoint: "",
				Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
				Children: []*Node{
					{
						Depth:    1,
						Endpoint: "api",
						Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetAPI"}},
						Children: []*Node{
							{
								Depth:    2,
								Endpoint: "users",
								Methods: map[string]HandlerFunc{
									http.MethodGet:  {"handler", "GetUsers"},
									http.MethodPost: {"handler", "CreateUsers"},
								},
								Children: nil,
							},
							{
								Depth:    2,
								Endpoint: "products",
								Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetProducts"}},
							},
						},
					},
				},
			},
			args: args{
				p:           "/api/users/create",
				httpMethod:  http.MethodPost,
				handlerFunc: HandlerFunc{Package: "handler", Func: "CreateUsers"},
			},
			wantErr: false,
			wantField: fields{
				Depth:    0,
				Endpoint: "",
				Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
				Children: []*Node{
					{
						Depth:    1,
						Endpoint: "api",
						Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetAPI"}},
						Children: []*Node{
							{
								Depth:    2,
								Endpoint: "users",
								Methods: map[string]HandlerFunc{
									http.MethodGet:  {"handler", "GetUsers"},
									http.MethodPost: {"handler", "CreateUsers"},
								},
								Children: []*Node{
									{
										Depth:    3,
										Endpoint: "create",
										Methods:  map[string]HandlerFunc{http.MethodPost: {"handler", "CreateUsers"}},
										Children: nil,
									},
								},
							},
							{
								Depth:    2,
								Endpoint: "products",
								Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetProducts"}},
							},
						},
					},
				},
			},
		},
		{
			name: "Add path parameter to users",
			fields: fields{
				Depth:    0,
				Endpoint: "",
				Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
				Children: []*Node{
					{
						Depth:    1,
						Endpoint: "api",
						Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetAPI"}},
						Children: []*Node{
							{
								Depth:    2,
								Endpoint: "users",
								Methods: map[string]HandlerFunc{
									http.MethodGet:  {"handler", "GetUsers"},
									http.MethodPost: {"handler", "CreateUsers"},
								},
								Children: []*Node{
									{
										Depth:    3,
										Endpoint: "create",
										Methods:  map[string]HandlerFunc{http.MethodPost: {"handler", "CreateUsers"}},
										Children: nil,
									},
								},
							},
							{
								Depth:    2,
								Endpoint: "products",
								Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetProducts"}},
							},
						},
					},
				},
			},
			args: args{
				p:           "/api/users/:user_id",
				httpMethod:  http.MethodGet,
				handlerFunc: HandlerFunc{Package: "handler", Func: "GetUser"},
			},
			wantErr: false,
			wantField: fields{
				Depth:    0,
				Endpoint: "",
				Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
				Children: []*Node{
					{
						Depth:    1,
						Endpoint: "api",
						Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetAPI"}},
						Children: []*Node{
							{
								Depth:    2,
								Endpoint: "users",
								Methods: map[string]HandlerFunc{
									http.MethodGet:  {"handler", "GetUsers"},
									http.MethodPost: {"handler", "CreateUsers"},
								},
								Children: []*Node{
									{
										Depth:    3,
										Endpoint: "create",
										Methods:  map[string]HandlerFunc{http.MethodPost: {"handler", "CreateUsers"}},
										Children: nil,
									},
									{
										Depth:       3,
										Endpoint:    "user_id",
										IsPathParam: true,
										Methods:     map[string]HandlerFunc{http.MethodGet: {"handler", "GetUser"}},
									},
								},
							},
							{
								Depth:    2,
								Endpoint: "products",
								Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetProducts"}},
							},
						},
					},
				},
			},
		},
		{
			name: "Add [GET] posts",
			fields: fields{
				Depth:    0,
				Endpoint: "",
				Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
				Children: []*Node{
					{
						Depth:    1,
						Endpoint: "api",
						Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetAPI"}},
						Children: []*Node{
							{
								Depth:    2,
								Endpoint: "users",
								Methods: map[string]HandlerFunc{
									http.MethodGet:  {"handler", "GetUsers"},
									http.MethodPost: {"handler", "CreateUsers"},
								},
								Children: []*Node{
									{
										Depth:    3,
										Endpoint: "create",
										Methods:  map[string]HandlerFunc{http.MethodPost: {"handler", "CreateUsers"}},
										Children: nil,
									},
									{
										Depth:       3,
										Endpoint:    "user_id",
										IsPathParam: true,
										Methods:     map[string]HandlerFunc{http.MethodGet: {"handler", "GetUser"}},
									},
								},
							},
							{
								Depth:    2,
								Endpoint: "products",
								Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetProducts"}},
							},
						},
					},
				},
			},
			args: args{
				p:           "/api/users/:user_id/posts",
				httpMethod:  http.MethodGet,
				handlerFunc: HandlerFunc{Package: "handler", Func: "GetPosts"},
			},
			wantErr: false,
			wantField: fields{
				Depth:    0,
				Endpoint: "",
				Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
				Children: []*Node{
					{
						Depth:    1,
						Endpoint: "api",
						Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetAPI"}},
						Children: []*Node{
							{
								Depth:    2,
								Endpoint: "users",
								Methods: map[string]HandlerFunc{
									http.MethodGet:  {"handler", "GetUsers"},
									http.MethodPost: {"handler", "CreateUsers"},
								},
								Children: []*Node{
									{
										Depth:    3,
										Endpoint: "create",
										Methods:  map[string]HandlerFunc{http.MethodPost: {"handler", "CreateUsers"}},
										Children: nil,
									},
									{
										Depth:       3,
										Endpoint:    "user_id",
										IsPathParam: true,
										Methods:     map[string]HandlerFunc{http.MethodGet: {"handler", "GetUser"}},
										Children: []*Node{
											{
												Depth:    4,
												Endpoint: "posts",
												Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetPosts"}},
												Children: nil,
											},
										},
									},
								},
							},
							{
								Depth:    2,
								Endpoint: "products",
								Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetProducts"}},
							},
						},
					},
				},
			},
		},
		{
			name: "Add [GET] path parameter to posts",
			fields: fields{
				Depth:    0,
				Endpoint: "",
				Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
				Children: []*Node{
					{
						Depth:    1,
						Endpoint: "api",
						Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetAPI"}},
						Children: []*Node{
							{
								Depth:    2,
								Endpoint: "users",
								Methods: map[string]HandlerFunc{
									http.MethodGet:  {"handler", "GetUsers"},
									http.MethodPost: {"handler", "CreateUsers"},
								},
								Children: []*Node{
									{
										Depth:    3,
										Endpoint: "create",
										Methods:  map[string]HandlerFunc{http.MethodPost: {"handler", "CreateUsers"}},
										Children: nil,
									},
									{
										Depth:       3,
										Endpoint:    "user_id",
										IsPathParam: true,
										Methods:     map[string]HandlerFunc{http.MethodGet: {"handler", "GetUser"}},
										Children: []*Node{
											{
												Depth:    4,
												Endpoint: "posts",
												Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetPosts"}},
												Children: nil,
											},
										},
									},
								},
							},
							{
								Depth:    2,
								Endpoint: "products",
								Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetProducts"}},
							},
						},
					},
				},
			},
			args: args{
				p:           "/api/users/:user_id/posts/:post_id",
				httpMethod:  http.MethodGet,
				handlerFunc: HandlerFunc{Package: "handler", Func: "GetPost"},
			},
			wantErr: false,
			wantField: fields{
				Depth:    0,
				Endpoint: "",
				Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
				Children: []*Node{
					{
						Depth:    1,
						Endpoint: "api",
						Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetAPI"}},
						Children: []*Node{
							{
								Depth:    2,
								Endpoint: "users",
								Methods: map[string]HandlerFunc{
									http.MethodGet:  {"handler", "GetUsers"},
									http.MethodPost: {"handler", "CreateUsers"},
								},
								Children: []*Node{
									{
										Depth:    3,
										Endpoint: "create",
										Methods:  map[string]HandlerFunc{http.MethodPost: {"handler", "CreateUsers"}},
										Children: nil,
									},
									{
										Depth:       3,
										Endpoint:    "user_id",
										IsPathParam: true,
										Methods:     map[string]HandlerFunc{http.MethodGet: {"handler", "GetUser"}},
										Children: []*Node{
											{
												Depth:    4,
												Endpoint: "posts",
												Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetPosts"}},
												Children: []*Node{
													{
														Depth:       5,
														Endpoint:    "post_id",
														IsPathParam: true,
														Methods:     map[string]HandlerFunc{http.MethodGet: {"handler", "GetPost"}},
													},
												},
											},
										},
									},
								},
							},
							{
								Depth:    2,
								Endpoint: "products",
								Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetProducts"}},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pn := &Node{
				Endpoint:    tt.fields.Endpoint,
				IsPathParam: tt.fields.IsPathParam,
				Methods:     tt.fields.Methods,
				Parent:      tt.fields.Parent,
				Children:    tt.fields.Children,
			}
			pn.registerParent()
			if err := pn.Add(tt.args.p, tt.args.httpMethod, tt.args.handlerFunc); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			wantPn := &Node{
				Endpoint:    tt.wantField.Endpoint,
				IsPathParam: tt.wantField.IsPathParam,
				Methods:     tt.wantField.Methods,
				Parent:      tt.wantField.Parent,
				Children:    tt.wantField.Children,
			}
			wantPn.registerParent()
			if !reflect.DeepEqual(pn, wantPn) {
				fmt.Println("got:")
				pn.Print()
				fmt.Println("want:")
				wantPn.Print()
				t.Errorf("Add() got = %v, want %v", pn, wantPn)
			}
		})
	}
}

func TestWalk(t *testing.T) {
	//TODO add parent
	endPoints := make([]string, 0)
	buf := &bytes.Buffer{}
	node := &Node{
		Depth:    0,
		Endpoint: "",
		Methods:  map[string]HandlerFunc{"GET": {"handler", "RootGet"}},
		Children: []*Node{
			{
				Depth:    1,
				Endpoint: "api",
				Methods:  map[string]HandlerFunc{"GET": {"handler", "API"}},
				Children: []*Node{
					{
						Depth:    2,
						Endpoint: "products",
						Methods:  map[string]HandlerFunc{"GET": {"handler", "GetProducts"}},
					},
					{
						Depth:    2,
						Endpoint: "users",
						Methods: map[string]HandlerFunc{
							"GET":  {"handler", "GetUsers"},
							"POST": {"handler", "CreateUser"},
						},
						Children: []*Node{
							{
								Depth:       3,
								Endpoint:    "user_id",
								IsPathParam: true,
								Methods:     map[string]HandlerFunc{"GET": {"handler", "GetUser"}},
								Children: []*Node{
									{
										Depth:    4,
										Endpoint: "posts",
										Methods:  map[string]HandlerFunc{"GET": {"handler", "GetPosts"}},
										Children: []*Node{
											{
												Depth:       5,
												Endpoint:    "post_id",
												IsPathParam: true,
												Methods:     map[string]HandlerFunc{"GET": {"handler", "GetPost"}},
												Children:    nil,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	type args struct {
		node *Node
		f    func(*Node) bool
	}
	tests := []struct {
		name          string
		args          args
		wantEndpoints []string
		wantOutput    string
	}{
		{
			name: "Print info",
			args: args{
				node: node,
				f: func(node *Node) bool {
					_, err := fmt.Fprintf(buf, "Endpoint: %v, IsPathPram: %v, Methods: %v\n", node.Endpoint, node.IsPathParam, node.Methods)
					if err != nil {
						log.Fatalln(err)
					}
					endPoints = append(endPoints, node.Endpoint)
					return true
				},
			},
			wantEndpoints: []string{"", "api", "products", "users", "user_id", "posts", "post_id"},
			wantOutput: `Endpoint: , IsPathPram: false, Methods: map[GET:{handler RootGet}]
Endpoint: api, IsPathPram: false, Methods: map[GET:{handler API}]
Endpoint: products, IsPathPram: false, Methods: map[GET:{handler GetProducts}]
Endpoint: users, IsPathPram: false, Methods: map[GET:{handler GetUsers} POST:{handler CreateUser}]
Endpoint: user_id, IsPathPram: true, Methods: map[GET:{handler GetUser}]
Endpoint: posts, IsPathPram: false, Methods: map[GET:{handler GetPosts}]
Endpoint: post_id, IsPathPram: true, Methods: map[GET:{handler GetPost}]
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Walk(tt.args.node, tt.args.f)
			for i, v := range endPoints {
				if v != tt.wantEndpoints[i] {
					t.Errorf("Unexpected output: \ngot:\n%s \nwant: \n%s", v, tt.wantEndpoints[i])
				}
			}
			if buf.String() != tt.wantOutput {
				t.Errorf("Unexpected output: \ngot:\n%s \nwant: \n%s", buf.String(), tt.wantOutput)
			}
		})
	}
}

func TestBuildBasePath(t *testing.T) {
	rootNode := &Node{
		Depth:    0,
		Endpoint: "",
		Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
		Children: []*Node{
			{
				Depth:    1,
				Endpoint: "api",
				Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetAPI"}},
				Children: []*Node{
					{
						Depth:    2,
						Endpoint: "users",
						Methods: map[string]HandlerFunc{
							http.MethodGet:  {"handler", "GetUsers"},
							http.MethodPost: {"handler", "CreateUsers"},
						},
					},
				},
			},
		},
	}
	node := &Node{
		Depth:    3,
		Endpoint: "create",
		Methods:  map[string]HandlerFunc{http.MethodPost: {"handler", "CreateUser"}},
		Children: nil,
	}
	rootNode.Children[0].Children[0].Children = append(rootNode.Children[0].Children[0].Children, node)
	rootNode.registerParent()

	rootNode2 := &Node{
		Depth:    0,
		Endpoint: "",
		Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
		Children: []*Node{
			{
				Depth:    1,
				Endpoint: "api",
				Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetAPI"}},
				Children: []*Node{
					{
						Depth:    2,
						Endpoint: "users",
						Methods: map[string]HandlerFunc{
							http.MethodGet:  {"handler", "GetUsers"},
							http.MethodPost: {"handler", "CreateUsers"},
						},
						Children: []*Node{
							{
								Depth:    3,
								Endpoint: "create",
								Methods:  map[string]HandlerFunc{http.MethodPost: {"handler", "CreateUsers"}},
							},
							{
								Depth:       3,
								Endpoint:    "user_id",
								IsPathParam: true,
								Methods:     map[string]HandlerFunc{http.MethodGet: {"handler", "GetUser"}},
								Children: []*Node{
									{
										Depth:    4,
										Endpoint: "posts",
										Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetPosts"}},
										Children: nil,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	node2 := &Node{
		Depth:    5,
		Endpoint: "create",
		Methods:  map[string]HandlerFunc{http.MethodPost: {"handler", "CreatePost"}},
		Children: nil,
	}
	rootNode2.Children[0].Children[0].Children[1].Children[0].Children = append(rootNode2.Children[0].Children[0].Children[1].Children[0].Children, node2)
	rootNode2.registerParent()

	rootNode3 := &Node{
		Depth:    0,
		Endpoint: "",
		Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetRoot"}},
		Children: nil,
	}
	node3 := &Node{
		Depth:1,
		Endpoint:"api",
		Methods:  map[string]HandlerFunc{http.MethodGet: {"handler", "GetAPI"}},
	}
	rootNode3.Children = append(rootNode3.Children, node3)
	rootNode3.registerParent()
	type args struct {
		node *Node
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "if path param does not exist, build base path from arg node to root node",
			args: args{node: node},
			want: "/api/users",
		},
		{
			name: "if path param exists, build base path from arg node to path param node",
			args: args{node: node2},
			want: "/posts",
		},
		{
			name: `if the depth of arg node is 1, return "/"`,
			args: args{node: node3},
			want: "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildBasePath(tt.args.node); got != tt.want {
				t.Errorf("BuildBasePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
