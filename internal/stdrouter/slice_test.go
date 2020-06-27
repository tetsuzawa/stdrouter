package stdrouter

import (
	"reflect"
	"testing"
)

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
