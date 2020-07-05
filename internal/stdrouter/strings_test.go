package stdrouter

import "testing"

func TestSnakeToCamel(t *testing.T) {
	type args struct {
		snake_case string
	}
	tests := []struct {
		name          string
		args          args
		wantCamelCase string
	}{
		{
			name:          "convert snake_case to CamelCase",
			args:          args{"user_id"},
			wantCamelCase: "UserId",
		},
		{
			name:          "do nothing if the argument is CamelCase",
			args:          args{"UserId"},
			wantCamelCase: "UserId",
		},
		{
			name:          "if the first letter is lowercase, capitalize",
			args:          args{"userId"},
			wantCamelCase: "UserId",
		},
		{
			name:          "if the argument is mix of snake_case and CamelCase, make it proper CamelCase",
			args:          args{"User_Id"},
			wantCamelCase: "UserId",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCamelCase := SnakeToCamel(tt.args.snake_case); gotCamelCase != tt.wantCamelCase {
				t.Errorf("SnakeToCamel() = %v, want %v", gotCamelCase, tt.wantCamelCase)
			}
		})
	}
}

func TestToLowerFirstLetter(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "convert first letter to lowercase",
			args: args{"UserId"},
			want: "userId",
		},
		{
			name: "do nothing if the first letter of string is already lowercase",
			args: args{"userId"},
			want: "userId",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToLowerFirstLetter(tt.args.s); got != tt.want {
				t.Errorf("ToLowerFirstLetter() = %v, want %v", got, tt.want)
			}
		})
	}
}
