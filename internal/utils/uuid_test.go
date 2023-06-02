package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUUID(t *testing.T) {
	ok := CheckUUID(NewUUID())
	require.True(t, ok)

	nok := CheckUUID("12121212121212121212")
	require.False(t, nok)
}
