package safe_math

import (
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestSafeAdd(t *testing.T) {
	n := uint64(math.MaxUint64 - 1)
	res, err := SafeAdd(n, 1)
	require.NoError(t, err)
	require.Equal(t, res, uint64(math.MaxUint64))

	// overflow
	res, err = SafeAdd(n, 2)
	require.Error(t, err)
	require.Equal(t, res, uint64(0))
}
