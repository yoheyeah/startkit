package check

func CheckExt(ext string, exts []string) bool {
	for _, v := range exts {
		if ext == v {
			return true
		}
	}
	return false
}
