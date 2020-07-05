package stdrouter

import "testing"

func TestSeparatePath(t *testing.T) {
	const n = 2
	type args struct {
		p string
		n int
	}
	tests := []struct {
		name     string
		args     args
		wantHead string
		wantTail string
	}{
		{
			name:     "/",
			args:     args{p: "/", n: 0},
			wantHead: "/",
			wantTail: "",
		},
		{
			name:     "/",
			args:     args{p: "/", n: n},
			wantHead: "/",
			wantTail: "",
		},
		{
			name:     "/api",
			args:     args{p: "/api", n: n},
			wantHead: "/api",
			wantTail: "",
		},
		{
			name:     "/api/users",
			args:     args{p: "/api/users", n: n},
			wantHead: "/api/users",
			wantTail: "/",
		},
		{
			name:     "/api/users/create",
			args:     args{p: "/api/users/create", n: n},
			wantHead: "/api/users",
			wantTail: "/create",
		},
		{
			name:     "/api/users/1",
			args:     args{p: "/api/users/1", n: n},
			wantHead: "/api/users",
			wantTail: "/1",
		},
		{
			name:     "/api/users/1/posts",
			args:     args{p: "/api/users/1/posts", n: n},
			wantHead: "/api/users",
			wantTail: "/1/posts",
		},
		{
			name:     "/api/users/1/posts/5",
			args:     args{p: "/api/users/1/posts/5", n: n},
			wantHead: "/api/users",
			wantTail: "/1/posts/5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHead, gotTail := SeparatePath(tt.args.p, tt.args.n)
			if gotHead != tt.wantHead {
				t.Errorf("SeparatePath() gotHead = %v, want %v", gotHead, tt.wantHead)
			}
			if gotTail != tt.wantTail {
				t.Errorf("SeparatePath() gotTail = %v, want %v", gotTail, tt.wantTail)
			}
		})
	}
}
