package helper

import (
	"crypto/sha256"
	"fmt"
)

func Sha256Digest(pass string) string {
	true_hash := fmt.Sprintf("%x", sha256.Sum256([]byte(pass)))
	return true_hash
}
