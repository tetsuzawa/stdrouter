package stdrouter

import (
	"path"
	"strings"
)

// SeparatePath separate path with at the nth "/", then returns head and tail of path.
func SeparatePath(p string, n int) (head, tail string) {
	p = path.Clean("/" + p)
	ps := strings.Split(p[1:], "/")
	if len(ps) < 2 {
		return p, ""
	}
	head = path.Clean("/" + strings.Join(ps[:n], "/"))
	tail = path.Clean("/" + strings.Join(ps[n:], "/"))
	return head, tail
}
