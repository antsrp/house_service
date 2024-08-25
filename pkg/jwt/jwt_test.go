package jwt_test

import (
	"testing"

	"github.com/antsrp/house_service/pkg/jwt"
	"github.com/stretchr/testify/require"
)

func Test1(t *testing.T) {
	key := `some key`
	service := jwt.NewJwtService([]byte(key))

	fields := map[string]any{
		"token": "32354",
	}

	token, err := service.NewToken(fields)
	require.NoError(t, err)

	parsed, err := service.Parse(token)
	require.NoError(t, err)

	require.Equalf(t, len(fields), len(parsed), "expected count of fields: %d, real: %d", len(fields), len(parsed))
	for k, v := range fields {
		require.Equalf(t, v, parsed[k], "value of key %s should be %v, actual %v", k, v, parsed[k])
	}
}
