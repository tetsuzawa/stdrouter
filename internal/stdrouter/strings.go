package stdrouter

import (
	"strings"
	"unicode"
)

// SnakeToCamel converts string from snake_case to CamelCase.
func SnakeToCamel(snake_case string) (CamelCase string) {
	isToUpper := false
	for k, v := range snake_case {
		if k == 0 {
			CamelCase = strings.ToUpper(string(snake_case[0]))
		} else {
			if isToUpper {
				CamelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					CamelCase += string(v)
				}
			}
		}
	}
	return
}

//ToLowerFirstLetter converts ExampleString to exampleString.
func ToLowerFirstLetter(s string) string {
	a := []rune(s)
	a[0] = unicode.ToLower(a[0])
	s = string(a)
	return s
}
