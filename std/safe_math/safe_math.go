package safe_math

import (
	"errors"
	"strconv"
)

func SafeAdd(a, b uint64) (uint64, error) {
	res := a + b
	if res >= a && res >= b {
		return res, nil
	}
	return 0, errors.New("overflow in add")
}

func SafeSub(a, b uint64) (uint64, error) {
	if b > a {
		return 0, errors.New("large subtractor with" + strconv.Itoa(int(b)))
	}
	return a - b, nil
}

func SafeMul(a, b uint64) (uint64, error) {
	res := a * b
	if a == 0 || res/a == b {
		return res, nil
	}
	return 0, errors.New("overflow in mul")
}

func SafeDiv(a, b uint64) (uint64, error) {
	if b == 0 {
		return 0, errors.New("invalid divisor with 0")
	}
	res := a / b
	if a == b*res+a%b {
		return res, nil
	}
	return 0, errors.New("overflow in div")
}
