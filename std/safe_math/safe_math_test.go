package safe_math

import (
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestSafeAdd(t *testing.T) {
	specs := map[string]struct {
		a      uint64
		b      uint64
		expRes uint64
		expErr bool
	}{
		"pass_1": {
			uint64(math.MaxUint64 - 1),
			1,
			uint64(math.MaxUint64),
			false,
		},
		"pass_2": {
			1,
			uint64(math.MaxUint64 - 1),
			uint64(math.MaxUint64),
			false,
		},
		"overflow_1": {
			uint64(math.MaxUint64 - 1),
			2,
			0,
			true,
		},
		"overflow_2": {
			2,
			uint64(math.MaxUint64 - 1),
			0,
			true,
		},
	}

	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			res, err := SafeAdd(spec.a, spec.b)
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, spec.expRes, res)
		})
	}
}

func TestSafeSub(t *testing.T) {
	specs := map[string]struct {
		a      uint64
		b      uint64
		expRes uint64
		expErr bool
	}{
		"pass_1": {
			2,
			1,
			1,
			false,
		},
		"pass_2": {
			1,
			0,
			1,
			false,
		},
		"pass_3": {
			1,
			1,
			0,
			false,
		},
		"overflow": {
			0,
			1,
			0,
			true,
		},
	}

	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			res, err := SafeSub(spec.a, spec.b)
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, spec.expRes, res)
		})
	}
}

func TestSafeMul(t *testing.T) {
	n := uint64(math.MaxUint64 / 5)
	res, err := SafeMul(n, 5)
	require.NoError(t, err)
	require.Equal(t, res, uint64(math.MaxUint64))

	res, err = SafeMul(0, n)
	require.NoError(t, err)
	require.Equal(t, res, uint64(0))

	// overflow
	res, err = SafeMul(n, 6)
	require.Error(t, err)
	require.Equal(t, res, uint64(0))
}

func TestSafeDiv(t *testing.T) {
	n := uint64(math.MaxUint64 / 5)
	res, err := SafeDiv(uint64(math.MaxUint64), 5)
	require.NoError(t, err)
	require.Equal(t, res, n)

	// overflow
	res, err = SafeDiv(uint64(math.MaxUint64), 0)
	require.Error(t, err)
	require.Equal(t, res, uint64(0))
}
