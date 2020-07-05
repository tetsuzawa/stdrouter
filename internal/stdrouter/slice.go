package stdrouter

// DropDuplication drops duplicate elements.
func DropDuplication(ss []string) []string {
	if ss == nil {
		return nil
	}
	res := make([]string, 0, len(ss))
	encountered := map[string]bool{}
	for i := 0; i < len(ss); i++ {
		if !encountered[ss[i]] {
			encountered[ss[i]] = true
			res = append(res, ss[i])
		}
	}
	return res
}

// Contains checks whether the argument slice s contains e.
func Contains(s string, ss []string) bool {
	for _, v := range ss {
		if s == v {
			return true
		}
	}
	return false
}
