package hash

import (
	"crypto/sha1"
	"fmt"
)

// salt = hex
type SHA1Hasher struct {
	Salt string
}

func NewSHA1Hasher(salt string) *SHA1Hasher {
	return &SHA1Hasher{Salt: salt}
}

func (h *SHA1Hasher) Hash(password string) (string, error) {
	hash := sha1.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return "", nil
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.Salt))), nil
}
