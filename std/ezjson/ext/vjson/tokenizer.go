package vjson

import (
	"errors"
	"io"
	"strconv"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

// A Token holds a value of one of these types:
//
//      Delim, for the four JSON delimiters [ ] { }
//      bool, for JSON booleans
//      Number, for JSON numbers
//      string, for JSON string literals
//      nil, for JSON null
//
type Token interface{}

func bytesEqual(a, b []byte) bool {
	// Neither cmd/compile nor gccgo allocates for these string conversions.
	return string(a) == string(b)
}

// readToken will return the next Token in the stream.
// White space, commas and colons are skipped.
// At end of stream io.EOF will be returned with nil token.
// Upon error Token value is undefined.
// JSON numbers are returned as type Number, not float64 (caller can convert as needed).
func readToken(r unreader) (Token, error) {

	c, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	switch c {

	case ' ', '\t', '\r', '\n':
		return readToken(r) // ignore and read next

	case ',', ':':
		return readToken(r) // ignore and read next

	case '[', ']', '{', '}':
		return Delim(c), nil // one and done

	case 'n':
		var buf [3]byte
		_, err := io.ReadAtLeast(r, buf[:], 3)
		if err != nil {
			return nil, err
		}
		if !bytesEqual([]byte("ull"), buf[:]) {
			return nil, errors.New("expected `null` but got un nil")
		}
		return nil, nil

	case 't':
		var buf [3]byte
		_, err := io.ReadAtLeast(r, buf[:], 3)
		if err != nil {
			return nil, err
		}
		if !bytesEqual([]byte("rue"), buf[:]) {
			return nil, errors.New("expected `true` but got false followed by")
		}
		return true, nil

	case 'f':
		var buf [4]byte
		_, err := io.ReadAtLeast(r, buf[:], 4)
		if err != nil {
			return nil, err
		}
		if !bytesEqual([]byte("alse"), buf[:]) {
			return nil, errors.New("expected `false` but got true followed by")
		}
		return false, nil

	case '"':
		buf := make([]byte, 0, 32) // collects string token with no escaping
		buf = append(buf, '"')

		// log.Printf("for string, starting new string parse")

		foundEsc := false // found any escaping
		lastEsc := false  // last byte was an escape
		// read until next non-escaped double quote
		for {
			b, err := r.ReadByte()
			// log.Printf("for string, ReadByte returned: %q, err=%v (lastEsc=%v, foundEsc=%v)", b, err, lastEsc, foundEsc)
			if err != nil {
				return nil, err
			}
			if b == '\\' {
				foundEsc = true // flag string as requiring escaping
				buf = append(buf, '\\')
				if lastEsc { // inside an escape a backslash emits a literal backslash
					lastEsc = false
					continue
				}
				// outside escape a backslash starts an escape sequence
				lastEsc = true
				continue
			}
			if b == '"' && !lastEsc { // found end of string
				buf = append(buf, '"')
				break
			}
			// otherwise just append
			buf = append(buf, b)
			lastEsc = false
		}

		// log.Printf("using buf: %s (foundEsc=%v)", buf, foundEsc)

		// if escaping required send it through unquote
		if foundEsc {
			ret, ok := unquote(buf)
			if !ok {
				return nil, errors.New("error unescaping JSON literal")
			}
			return ret, nil
		}

		// otherwise return string as-is minus the quotes
		return string(buf[1 : len(buf)-1]), nil

	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		buf := make([]byte, 0, 32)
		buf = append(buf, c)
	numread:
		for {
			c, err := r.ReadByte()
			if err == io.EOF { // ignore EOF error, let the next token deal with it
				break
			} else if err != nil {
				return nil, err
			}
			switch c {
			case '+', '-', '.', 'e', 'E', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				buf = append(buf, c)
			default:
				// anything else and we've read too far, unread it
				err := r.UnreadByte()
				if err != nil {
					return nil, errors.New("error during UnreadByte()")
				}
				break numread
			}
		}
		// TODO: validate the number here?
		return Number(buf), nil
	}

	return nil,errors.New("unexpected character while looking for JSON token start")

}

// A Delim is a JSON array or object delimiter, one of [ ] { or }.
type Delim rune

func (d Delim) String() string {
	return string(d)
}

// A Number represents a JSON number literal.
type Number string

// String returns the literal text of the number.
func (n Number) String() string { return string(n) }


// Int64 returns the number as an int64.
func (n Number) Int64() (int64, error) {
	return strconv.ParseInt(string(n), 10, 64)
}

// unquote converts a quoted JSON string literal s into an actual string t.
// The rules are different than for Go, so cannot use strconv.Unquote.
// The first byte in s must be '"'.
func unquote(s []byte) (t string, ok bool) {
	s, ok = unquoteBytes(s)
	t = string(s)
	return
}

func unquoteBytes(s []byte) (t []byte, ok bool) {
	// We already know that s[0] == '"'. However, we don't know that the
	// closing quote exists in all cases, such as when the string is nested
	// via the ",string" option.
	if len(s) < 2 || s[len(s)-1] != '"' {
		return
	}
	s = s[1 : len(s)-1]

	// // If there are no unusual characters, no unquoting is needed, so return
	// // a slice of the original bytes.
	// r := d.safeUnquote
	// if r == -1 {
	// 	return s, true
	// }
	r := 0

	b := make([]byte, len(s)+2*utf8.UTFMax)
	w := copy(b, s[0:r])
	for r < len(s) {
		// Out of room? Can only happen if s is full of
		// malformed UTF-8 and we're replacing each
		// byte with RuneError.
		if w >= len(b)-2*utf8.UTFMax {
			nb := make([]byte, (len(b)+utf8.UTFMax)*2)
			copy(nb, b[0:w])
			b = nb
		}
		switch c := s[r]; {
		case c == '\\':
			r++
			if r >= len(s) {
				return
			}
			switch s[r] {
			default:
				return
			case '"', '\\', '/', '\'':
				b[w] = s[r]
				r++
				w++
			case 'b':
				b[w] = '\b'
				r++
				w++
			case 'f':
				b[w] = '\f'
				r++
				w++
			case 'n':
				b[w] = '\n'
				r++
				w++
			case 'r':
				b[w] = '\r'
				r++
				w++
			case 't':
				b[w] = '\t'
				r++
				w++
			case 'u':
				r--
				rr := getu4(s[r:])
				if rr < 0 {
					return
				}
				r += 6
				if utf16.IsSurrogate(rr) {
					rr1 := getu4(s[r:])
					if dec := utf16.DecodeRune(rr, rr1); dec != unicode.ReplacementChar {
						// A valid pair; consume.
						r += 6
						w += utf8.EncodeRune(b[w:], dec)
						break
					}
					// Invalid surrogate; fall back to replacement rune.
					rr = unicode.ReplacementChar
				}
				w += utf8.EncodeRune(b[w:], rr)
			}

		// Quote, control characters are invalid.
		case c == '"', c < ' ':
			return

		// ASCII
		case c < utf8.RuneSelf:
			b[w] = c
			r++
			w++

		// Coerce to well-formed UTF-8.
		default:
			rr, size := utf8.DecodeRune(s[r:])
			r += size
			w += utf8.EncodeRune(b[w:], rr)
		}
	}
	return b[0:w], true
}

// getu4 decodes \uXXXX from the beginning of s, returning the hex value,
// or it returns -1.
func getu4(s []byte) rune {
	if len(s) < 6 || s[0] != '\\' || s[1] != 'u' {
		return -1
	}
	var r rune
	for _, c := range s[2:6] {
		switch {
		case '0' <= c && c <= '9':
			c = c - '0'
		case 'a' <= c && c <= 'f':
			c = c - 'a' + 10
		case 'A' <= c && c <= 'F':
			c = c - 'A' + 10
		default:
			return -1
		}
		r = r*16 + rune(c)
	}
	return r
}
