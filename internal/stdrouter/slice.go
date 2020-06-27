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
