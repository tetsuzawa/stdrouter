package stdrouter

import (
	"reflect"
	"testing"
)

func TestDrop(t *testing.T) {
	type args struct {
		s  string
		ss []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "drop if slice contains designated element",
			args: args{"bbb", []string{"aaa", "bbb", "ccc", "ddd"}},
			want: []string{"aaa", "ccc", "ddd"},
		},
		{
			name: "do nothing if slice does not contains designated element",
			args: args{"bbb", []string{"aaa", "ccc", "ddd"}},
			want: []string{"aaa", "ccc", "ddd"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Drop(tt.args.s, tt.args.ss); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Drop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDropDuplication(t *testing.T) {
	type args struct {
		ss []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "drop duplicate element",
			args: args{[]string{"aaa", "aaa", "bbb", "ccc", "ccc", "ddd"}},
			want: []string{"aaa", "bbb", "ccc", "ddd"},
		},
		{
			name: "do nothing if slice is empty",
			args: args{[]string{}},
			want: []string{},
		},
		{
			name: "do nothing if slice is nil",
			args: args{nil},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DropDuplication(tt.args.ss); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DropDuplication() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args struct {
		s  string
		ss []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "return true if the slice contains",
			args: args{"abc", []string{"abc", "def", "ghi"}},
			want: true,
		},
		{
			name: "return false if the slice does not contains",
			args: args{"xyz", []string{"abc", "def", "ghi"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.s, tt.args.ss); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
