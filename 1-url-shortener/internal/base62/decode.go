package base62

import "fmt"

func DecodeBase62(base62str string) (int64, error) {
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
