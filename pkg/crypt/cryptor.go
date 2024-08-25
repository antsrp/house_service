package crypt

import (
	"crypto/md5"
	"encoding/hex"
)

type Cryptor interface {
	Hash(string) string
}

type Crypt struct {
}

var _ Cryptor = Crypt{}

func (crypt Crypt) Hash(s string) string {
	data := md5.Sum([]byte(s))
	return hex.EncodeToString(data[:])
}
