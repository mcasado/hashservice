package app

import (
	"crypto/sha512"
	"encoding/base64"
)

func CreateHash(input string) string {

	sha512 := sha512.New()
	sha512.Write([]byte(input))
	return base64.StdEncoding.EncodeToString(sha512.Sum(nil))
}

func Max(numbers []int) int {

	max := 0

	// Loop through all integers in array.
	for _, number := range numbers {
		// Check if number is great than max.
		if number > max {
			max = number
		}
	}
	return max
}
