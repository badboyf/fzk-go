package util

var Util = &util{}

type util struct {
}

func (u util) Contains(v string, vs []string) bool {
	if vs == nil {
		return false
	}
	for _, tmp := range vs {
		if v == tmp {
			return true
		}
	}
	return false
}