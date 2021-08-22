package std

import (
	"strconv"

	"github.com/cosmwasm/cosmwasm-go/std/types"
)

func SafeAdd(a, b uint64) (uint64, error) {
	res := a + b
	if res >= a && res >= b {
		return res, nil
	}
	return 0, types.NewError("overflow in add")
}

func SafeSub(a, b uint64) (uint64, error) {
	if b > a {
		return 0, types.NewError("large subtractor with" + strconv.Itoa(int(b)))
	}
	return a - b, nil
}

func SafeMul(a, b uint64) (uint64, error) {
	res := a * b
	if a == 0 || res/a == b {
		return res, nil
	}
	return 0, types.NewError("overflow in mul")
}

func SafeDiv(a, b uint64) (uint64, error) {
	if b == 0 {
		return 0, types.NewError("invalid divisor with 0")
	}
	res := a / b
	if a == b*res+a%b {
		return res, nil
	}
	return 0, types.NewError("overflow in div")
}
