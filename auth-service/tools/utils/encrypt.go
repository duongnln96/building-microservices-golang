package utils

import (
	"crypto/sha1"
	"fmt"
)

func Encrypt(s string) string {
	return fmt.Sprintf("%s", sha1.Sum([]byte(s)))
}
