package vjson

import (
	"io"
	"unicode/utf8"
)

var hex = "0123456789abcdef"
// NOTE: keep in sync with encodeStringBytes below.
func encodeString(w io.Writer, s string, escapeHTML bool) error {
	w.Write([]byte{'"'})
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] || (!escapeHTML && safeSet[b]) {
				i++
				continue
			}
			if start < i {
				w.Write([]byte(s[start:i]))
			}
			switch b {
			case '\\', '"':
				w.Write([]byte{b})
			case '\n':
				w.Write([]byte{'n'})
			case '\r':
				w.Write([]byte{'r'})
			case '\t':
				w.Write([]byte{'t'})
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				w.Write([]byte(`u00`))
				w.Write([]byte{hex[b>>4]})
				w.Write([]byte{hex[b&0xF]})
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				w.Write([]byte(s[start:i]))
			}
			w.Write([]byte(`\ufffd`))
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				w.Write([]byte(s[start:i]))
			}
			w.Write([]byte(`\u202`))
			w.Write([]byte{hex[c&0xF]})
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		w.Write([]byte(s[start:]))
	}
	_, err := w.Write([]byte{'"'})
	return err
}

// NOTE: keep in sync with encodeString above.
func encodeStringBytes(w io.Writer, s []byte, escapeHTML bool) error {
	w.Write([]byte{'"'})
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] || (!escapeHTML && safeSet[b]) {
				i++
				continue
			}
			if start < i {
				w.Write(s[start:i])
			}
			w.Write([]byte{'\\'})
			switch b {
			case '\\', '"':
				w.Write([]byte{b})
			case '\n':
				w.Write([]byte{'n'})
			case '\r':
				w.Write([]byte{'r'})
			case '\t':
				w.Write([]byte{'t'})
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				w.Write([]byte(`u00`))
				w.Write([]byte{hex[b>>4]})
				w.Write([]byte{hex[b&0xF]})
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRune(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				w.Write(s[start:i])
			}
			w.Write([]byte(`\ufffd`))
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				w.Write(s[start:i])
			}
			w.Write([]byte(`\u202`))
			w.Write([]byte{hex[c&0xF]})
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		w.Write(s[start:])
	}
	_, err := w.Write([]byte{'"'})
	return err
}
