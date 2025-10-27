package base62

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func reverseBytes(bytes []byte) []byte {
	left, right := 0, len(bytes)-1
	for left < right {
		bytes[left], bytes[right] = bytes[right], bytes[left]
		left++
		right--
	}
	return bytes
}

func getCharVal(char byte) int64 {
	if '0' <= char && char <= '9' {
		return int64(char - '0')
	} else if 'A' <= char && char <= 'Z' {
		return int64(char - 'A' + 10)
	} else if 'a' <= char && char <= 'z' {
		return int64(char - 'a' + 36)
	} else {
		return -1
	}
}
