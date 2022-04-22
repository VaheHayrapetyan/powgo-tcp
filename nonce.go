package gopow_tcp

import (
	"crypto/rand"
)

func newNonce() []byte {
	const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length := 30
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}

	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}

	return bytes
}
