package std

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type Wrapper struct {
	OptA EmptyStruct `json:"opt_a,omitempty"`
	OptB EmptyStruct `json:"opt_b,omitempty"`
}

func TestEmptyStructSerialization(t *testing.T) {
	empty := EmptyStruct{}
	bz, err := json.Marshal(empty)
	require.NoError(t, err)
	require.Equal(t, string(bz), "{}")

	var parsed EmptyStruct
	err = json.Unmarshal(bz, &parsed)
	require.NoError(t, err)
	require.Equal(t, parsed.Seen, true)
}

func TestWrapperUnmarshal(t *testing.T) {
	cases := map[string]struct {
		data string
		optA bool
		optB bool
	}{
		"both missing": {"{}", false, false},
		"opt a":        {`{"opt_a": {}}`, true, false},
		"opt b":        {`{"opt_b": {}}`, false, true},
		"both":         {`{"opt_a": {}, "opt_b": {}}`, true, true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			var wrap Wrapper
			err := json.Unmarshal([]byte(tc.data), &wrap)
			require.NoError(t, err)
			require.Equal(t, tc.optA, wrap.OptA.Seen)
			require.Equal(t, tc.optB, wrap.OptB.Seen)
		})
	}
}
