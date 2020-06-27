package stdrouter

import (
	"strings"
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
