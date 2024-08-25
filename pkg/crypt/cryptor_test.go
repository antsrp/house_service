package crypt_test

import (
	"testing"

	"github.com/antsrp/house_service/pkg/crypt"
	"github.com/stretchr/testify/assert"
)

var cryptor crypt.Cryptor

const (
	a = "aee1"
	b = "123as"
)

func TestMain(m *testing.M) {
	cryptor = crypt.Crypt{}
	m.Run()
}

func TestCrypt(t *testing.T) {
	tests := []struct {
		a        string
		b        string
		expected bool
	}{
		{a, b, false},
		{a, a, true},
		{b, b, true},
		{b, a, false},
	}

	for _, test := range tests {
		actual := cryptor.Hash(test.a) == cryptor.Hash(test.b)
		relation := "equal"
		if !test.expected {
			relation = "not " + relation
		}
		assert.Equalf(t, test.expected, actual, "hash of %s and %s should be %s", test.a, test.b, relation)
	}
}
