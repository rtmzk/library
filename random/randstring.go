package utils

import "crypto/rand"

const (
	Alphabet62 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	Alphabet32 = "abcdefghijklmnopqrstuvwxyz1234567890"
)

func RandString(letters string, n int) string {
	output := make([]byte, n)

	randomness := make([]byte, n)

	_, err := rand.Read(randomness)
	if err != nil {
		panic(err)
	}
	l := len(letters)

	for pos := range output {
		random := randomness[pos]

		randomPos := random % uint8(l)
		output[pos] = letters[randomPos]
	}

	return string(output)
}
