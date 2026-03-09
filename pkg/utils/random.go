package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateRandomString generates a random string of the specified length.
func GenerateRandomString(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

// GenerateDefaultUsername generates a default username in the format "USER_<random_string>".
func GenerateDefaultUsername() string {
	return fmt.Sprintf("USER_%s", GenerateRandomString(16))
}