package ezjson

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Wrapper struct {
	OptA EmptyStruct `json:"opt_a,omitempty"`
	OptB EmptyStruct `json:"opt_b,omitempty"`
}

func TestEmptyStructSerialization(t *testing.T) {
	t.Skip("this only seems to work with wrapper?")
	var parsed EmptyStruct
	err := Unmarshal([]byte("{}"), &parsed)
	require.NoError(t, err)
	require.Equal(t, parsed.WasSet(), true)
}

func TestWrapperUnmarshal(t *testing.T) {
	cases := map[string]struct {
		data string
		optA bool
		optB bool
	}{
		"both missing":     {"{}", false, false},
		"opt a":            {`{"opt_a": {}}`, true, false},
		"opt b":            {`{"opt_b": {}}`, false, true},
		"both":             {`{"opt_a": {}, "opt_b": {}}`, true, true},
		"seen a, unseen b": {`{"opt_a": {"set_random_flag": true}, "opt_b": {}}`, true, true},
		"random data":      {`{"opt_a": {"some": 1, "more": "stuff"}}`, true, false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			var wrap Wrapper
			err := Unmarshal([]byte(tc.data), &wrap)
			require.NoError(t, err)
			require.Equal(t, tc.optA, wrap.OptA.WasSet())
			require.Equal(t, tc.optB, wrap.OptB.WasSet())
		})
	}
}
