package base62

func IsValidBase62(code string) bool {
	if code == "" {
		return false
	}
	for _, char := range code {
		switch {
		case '0' <= char && char <= '9':
		case 'A' <= char && char <= 'Z':
		case 'a' <= char && char <= 'z':
		default:
			return false
		}
	}
	return true
}
