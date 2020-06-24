package src

import "testing"
import "../std"

func TestAllocate(t *testing.T) {
	result := std.Package_message([]byte("1234567"))
	_ = result
}