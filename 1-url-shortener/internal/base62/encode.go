package base62

func EncodeBase62(value int64) string {
	if value == 0 {
		return "0"
	}
	base62Bytes := make([]byte, 0)

	for value > 0 {
		remainder := value % 62
		value /= 62
		base62Bytes = append(base62Bytes, base62Chars[remainder])
	}
	return string(reverseBytes(base62Bytes))
}
