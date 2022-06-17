package kbjson

import (
	"fmt"
	"strings"
	"unicode/utf8"
	"unsafe"
)

const (
	PadRight int = iota //向右填充字符
	PadLeft             //向左填充字符
)

type (
	stringS struct {
		str unsafe.Pointer
		len int
	}
	sliceT struct {
		arr unsafe.Pointer
		len int
		cap int
	}
)

// Bytes2String bytes to string
func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// String2Bytes string to bytes
func String2Bytes(s string) []byte {
	var b []byte
	str := (*stringS)(unsafe.Pointer(&s))
	pbytes := (*sliceT)(unsafe.Pointer(&b))
	pbytes.arr = str.str
	pbytes.len = str.len
	pbytes.cap = str.len
	return b
}

// Len string length (utf8)
func Len(str string) int {
	return utf8.RuneCountInString(str)
}

// Pad String padding
func Pad(raw string, length int, padStr string, padType int) string {
	l := length - Len(raw)
	if l <= 0 {
		return raw
	}
	if padType == PadRight {
		raw = fmt.Sprintf("%s%s", raw, strings.Repeat(padStr, l))
	} else if padType == PadLeft {
		raw = fmt.Sprintf("%s%s", strings.Repeat(padStr, l), raw)
	} else {
		left := 0
		right := 0
		if l > 1 {
			left = l / 2
			right = (l / 2) + (l % 2)
		}

		raw = fmt.Sprintf("%s%s%s", strings.Repeat(padStr, left), raw, strings.Repeat(padStr, right))
	}
	return raw
}

// TrimSpace TrimSpace
func TrimSpace(s string) string {
	space := [...]uint8{127, 128, 133, 160, 194, 226, 227}
	well := func(s uint8) bool {
		for i := range space {
			if space[i] == s {
				return true
			}
		}
		return false
	}

	for len(s) > 0 {
		if (s[0] <= 31) || s[0] <= ' ' || well(s[0]) {
			s = s[1:]
			continue
		}
		break
	}

	for len(s) > 0 {
		if s[len(s)-1] <= ' ' || (s[len(s)-1] <= 31) || well(s[len(s)-1]) {
			s = s[:len(s)-1]
			continue
		}
		break
	}
	return s
}

//
//
//--------
//字符串匹配
//字符串匹配
func Match(str, pattern string) bool {
	if pattern == "*" {
		return true
	}
	return deepMatch(str, pattern)
}

func deepMatch(str, pattern string) bool {
	for len(pattern) > 0 {
		if pattern[0] > 0x7f {
			return deepMatchRune(str, pattern)
		}
		switch pattern[0] {
		default:
			if len(str) == 0 {
				return false
			}
			if str[0] > 0x7f {
				return deepMatchRune(str, pattern)
			}
			if str[0] != pattern[0] {
				return false
			}
		case '?':
			if len(str) == 0 {
				return false
			}
		case '*':
			return deepMatch(str, pattern[1:]) ||
				(len(str) > 0 && deepMatch(str[1:], pattern))
		}
		str = str[1:]
		pattern = pattern[1:]
	}
	return len(str) == 0 && len(pattern) == 0
}

func deepMatchRune(str, pattern string) bool {
	var sr, pr rune
	var srsz, prsz int

	x7f := func(isStr bool) (r rune, p int) {
		var s uint8
		if isStr {
			s = str[0]
		} else {
			s = pattern[0]
		}
		if str[0] > 0x7f {
			r, p = utf8.DecodeRuneInString(str)
		} else {
			r, p = rune(s), 1
		}
		return
	}

	if len(str) > 0 {
		sr, srsz = x7f(true)
	} else {
		sr, srsz = utf8.RuneError, 0
	}
	if len(pattern) > 0 {
		pr, prsz = x7f(false)
	} else {
		pr, prsz = utf8.RuneError, 0
	}
	for pr != utf8.RuneError {
		switch pr {
		default:
			if srsz == utf8.RuneError {
				return false
			}
			if sr != pr {
				return false
			}
		case '?':
			if srsz == utf8.RuneError {
				return false
			}
		case '*':
			return deepMatchRune(str, pattern[prsz:]) ||
				(srsz > 0 && deepMatchRune(str[srsz:], pattern))
		}
		str = str[srsz:]
		pattern = pattern[prsz:]
		if len(str) > 0 {
			sr, srsz = x7f(true)
		} else {
			sr, srsz = utf8.RuneError, 0
		}
		if len(pattern) > 0 {
			pr, prsz = x7f(false)
		} else {
			pr, prsz = utf8.RuneError, 0
		}
	}
	return srsz == 0 && prsz == 0
}
