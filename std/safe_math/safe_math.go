package safe_math

import (
	"errors"
	"strconv"
)

func SafeAdd(a, b uint64) (res uint64, err error) {
	res = a + b
	if res >= a && res >= b {
		return
	}
	return res, errors.New("overflow in add")
}

func SafeSub(a, b uint64) (res uint64, err error) {
	if b > a {
		return res, errors.New("large subtractor with" + strconv.Itoa(int(b)))
	}
	return a - b, err
}

func SafeMul(a, b uint64) (res uint64, err error) {
	res = a * b
	if a == 0 || res/a == b {
		return
	}
	return 0, errors.New("overflow in mul")
}

func SafeDiv(a, b uint64) (res uint64, err error) {
	if b == 0 {
		return res, errors.New("invalid divisor with 0")
	}
	res = a / b
	if a == b*res+a%b {
		return
	}
	return res, errors.New("overflow in div")
}
