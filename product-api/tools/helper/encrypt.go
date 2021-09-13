package helper

import (
	"crypto/sha1"
	"fmt"
)

// hash plaintext with SHA-1
func Encrypt(plaintext string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
}
