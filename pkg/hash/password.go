package hash

import (
	"crypto/sha1"
)

type PasswordHasher interface {
	Hash(password string) ([]byte, error)
}

type SHA1Hasher struct {
	salt string
}

func NewSHA1Hasher(salt string) *SHA1Hasher {
	return &SHA1Hasher{salt: salt}
}

func (h *SHA1Hasher) Hash(password string) ([]byte, error) {
    hash := sha1.New()

    if _, err := hash.Write([]byte(password)); err != nil {
        return nil, err
    }

    return hash.Sum([]byte(h.salt)), nil
}
