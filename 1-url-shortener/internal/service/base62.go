package service

import "fmt"

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Encode(value int64) string {
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

func Decode(base62str string) (int64, error) {
	var value int64 = 0
	for _, c := range base62str {
		digit := getCharVal(byte(c))
		if digit == -1 {
			return -1, fmt.Errorf("error - invalid character found - %c", c)
		}
		value = value*62 + digit
	}
	return value, nil
}

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
