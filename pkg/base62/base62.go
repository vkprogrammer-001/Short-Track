package base62

import (
	"errors"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Encode(number int64) string {
	if number == 0 {
		return string(charset[0])
	}
	
	var sb strings.Builder
	length := int64(len(charset))
	
	for number > 0 {
		rem := number % length
		sb.WriteByte(charset[rem])
		number = number / length
	}
	
	// Reverse the string
	res := sb.String()
	runes := []rune(res)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func Decode(code string) (int64, error) {
	var number int64
	length := int64(len(charset))
	
	for _, r := range code {
		idx := strings.IndexRune(charset, r)
		if idx == -1 {
			return 0, errors.New("invalid character in base62 code")
		}
		number = number*length + int64(idx)
	}
	return number, nil
}