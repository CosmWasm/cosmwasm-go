package wjson

import (
	"bytes"
	"unicode/utf8"
)

const (
	caseMask     = ^byte(0x20)
	kelvin       = '\u212a'
	smallLongEss = '\u017f'
)

func foldFunc(s []byte) func(s, t []byte) bool {
	nonLetter := false
	special := false
	for _, b := range s {
		if b >= utf8.RuneSelf {
			return bytes.EqualFold
		}
		upper := b & caseMask
		if upper < 'A' || upper > 'Z' {
			nonLetter = true
		} else if upper == 'K' || upper == 'S' {
			special = true
		}
	}
	if special {
		return equalFoldRight
	}
	if nonLetter {
		return asciiEqualFold
	}
	return simpleLetterEqualFold
}

func simpleLetterEqualFold(s, t []byte) bool {
	if len(s) != len(t) {
		return false
	}
	for i, b := range s {
		if b&caseMask != t[i]&caseMask {
			return false
		}
	}
	return true
}

func asciiEqualFold(s, t []byte) bool {
	if len(s) != len(t) {
		return false
	}
	for i, sb := range s {
		tb := t[i]
		if sb == tb {
			continue
		}
		if ('a' <= sb && sb <= 'z') || ('A' <= sb && sb <= 'Z') {
			if sb&caseMask != tb&caseMask {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func equalFoldRight(s, t []byte) bool {
	for _, sb := range s {
		if len(t) == 0 {
			return false
		}
		tb := t[0]
		if tb < utf8.RuneSelf {
			if sb != tb {
				sbUpper := sb & caseMask
				if 'A' <= sbUpper && sbUpper <= 'Z' {
					if sbUpper != tb&caseMask {
						return false
					}
				} else {
					return false
				}
			}
			t = t[1:]
			continue
		}

		tr, size := utf8.DecodeRune(t)
		switch sb {
		case 's', 'S':
			if tr != smallLongEss {
				return false
			}
		case 'k', 'K':
			if tr != kelvin {
				return false
			}
		default:
			return false
		}
		t = t[size:]

	}
	if len(t) > 0 {
		return false
	}
	return true
}
