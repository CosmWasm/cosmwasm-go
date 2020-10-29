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
			math.MaxUint64 - 1,
			1,
			math.MaxUint64,
			false,
		},
		"pass_2": {
			1,
			math.MaxUint64 - 1,
			math.MaxUint64,
			false,
		},
		"overflow_1": {
			math.MaxUint64 - 1,
			2,
			0,
			true,
		},
		"overflow_2": {
			2,
			math.MaxUint64 - 1,
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
		"pass_4": {
			0,
			0,
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
	// note: multiplicand = math.MaxUint64 / 5
	const multiplicand = 3689348814741910323
	specs := map[string]struct {
		a      uint64
		b      uint64
		expRes uint64
		expErr bool
	}{
		"pass_1": {
			multiplicand,
			5,
			math.MaxUint64,
			false,
		},
		"pass_2": {
			multiplicand,
			0,
			0,
			false,
		},
		"pass_3": {
			0,
			multiplicand,
			0,
			false,
		},
		"pass_4": {
			0,
			0,
			0,
			false,
		},
		"overflow": {
			multiplicand,
			6,
			0,
			true,
		},
	}

	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			res, err := SafeMul(spec.a, spec.b)
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, spec.expRes, res)
		})
	}
}

func TestSafeDiv(t *testing.T) {
	// note: quotient = math.MaxUint64 / 5
	const quotient = 3689348814741910323
	specs := map[string]struct {
		a      uint64
		b      uint64
		expRes uint64
		expErr bool
	}{
		"pass_1": {
			math.MaxUint64,
			5,
			quotient,
			false,
		},
		"pass_2": {
			0,
			math.MaxUint64,
			0,
			false,
		},
		"pass_3": {
			3,
			3,
			1,
			false,
		},
		"pass_4": {
			4,
			3,
			1,
			false,
		},
		"pass_5": {
			5,
			3,
			1,
			false,
		},
		"overflow": {
			math.MaxUint64,
			0,
			0,
			true,
		},
	}

	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			res, err := SafeDiv(spec.a, spec.b)
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, spec.expRes, res)
		})
	}
}
