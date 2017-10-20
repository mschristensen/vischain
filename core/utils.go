package core

import (
	"crypto/sha256"
	"fmt"
)

func Sha256(bytes []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(bytes))
}
