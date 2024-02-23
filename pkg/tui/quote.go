package tui

import (
	"bytes"
	"strconv"
	"unicode/utf8"
)

func bleh() {
	strconv.Quote("")
}

// Quote quotes each argument and joins them with a space.
// If passed to /bin/sh, the resulting string will be split back into the
// original arguments.
func Quote(in string) string {
	out := bytes.Buffer{}

	// for _, r := range in {
	quote(in, '"', &out)

	// out.WriteString(thing[1 : len(thing)-1])
	// }

	return out.String()

	// var buf bytes.Buffer

	// for i, arg := range args {
	// 	if i != 0 {
	// 		buf.WriteByte(' ')
	// 	}

	// 	quote(arg, &buf)
	// }

	// return buf.String()
}

const (
	lowerhex = "0123456789abcdef"
	upperhex = "0123456789ABCDEF"
)

func quote(word string, quote byte, buf *bytes.Buffer) {
	var (
		ASCIIonly   = false
		graphicOnly = false
	)

	// We want to try to produce a "nice" output. As such, we will
	// backslash-escape most characters, but if we encounter a space, or if we
	// encounter an extra-special char (which doesn't work with
	// backslash-escaping) we switch over to quoting the whole word. We do this
	// with a space because it's typically easier for people to read multi-word
	// arguments when quoted with a space rather than with ugly backslashes
	// everywhere.
	// origLen := buf.Len()

	if len(word) == 0 {
		// oops, no content
		buf.WriteString("")

		return
	}

	cur := word

	for len(cur) > 0 {
		runeValue, width := utf8.DecodeRuneInString(cur)

		cur = cur[width:]

		if width == 1 && runeValue == utf8.RuneError {
			buf.WriteString(`\x`)
			buf.WriteByte(lowerhex[runeValue>>4])
			buf.WriteByte(lowerhex[runeValue&0xF])

			continue
		}

		buf.Write(appendEscapedRune(nil, runeValue, quote, ASCIIonly, graphicOnly))
	}

	return
}

func appendEscapedRune(buf []byte, r rune, quote byte, ASCIIonly, graphicOnly bool) []byte {
	// if r == rune(quote) || r == '\\' { // always backslashed
	if r == rune(quote) { // always backslashed
		buf = append(buf, '\\')
		buf = append(buf, byte(r))

		return buf
	}

	if ASCIIonly {
		if r < utf8.RuneSelf && strconv.IsPrint(r) {
			buf = append(buf, byte(r))

			return buf
		}
	} else if strconv.IsPrint(r) || graphicOnly && isInGraphicList(r) {
		return utf8.AppendRune(buf, r)
	}

	switch r {
	case '\a':
		buf = append(buf, `\a`...)
	case '\b':
		buf = append(buf, `\b`...)
	case '\f':
		buf = append(buf, `\f`...)
	case '\n':
		buf = append(buf, `\n`...)
	case '\r':
		buf = append(buf, `\r`...)
	case '\t':
		buf = append(buf, `\t`...)
	case '\v':
		buf = append(buf, `\v`...)
	default:
		switch {
		case r < ' ' || r == 0x7f:
			buf = append(buf, `\x`...)
			buf = append(buf, lowerhex[byte(r)>>4])
			buf = append(buf, lowerhex[byte(r)&0xF])
		case !utf8.ValidRune(r):
			r = 0xFFFD

			fallthrough
		case r < 0x10000:
			buf = append(buf, `\u`...)
			for s := 12; s >= 0; s -= 4 {
				buf = append(buf, lowerhex[r>>uint(s)&0xF])
			}
		default:
			buf = append(buf, `\U`...)
			for s := 28; s >= 0; s -= 4 {
				buf = append(buf, lowerhex[r>>uint(s)&0xF])
			}
		}
	}

	return buf
}

// isInGraphicList reports whether the rune is in the isGraphic list. This separation
// from IsGraphic allows quoteWith to avoid two calls to IsPrint.
// Should be called only if IsPrint fails.
func isInGraphicList(r rune) bool {
	// We know r must fit in 16 bits - see makeisprint.go.
	if r > 0xFFFF {
		return false
	}

	rr := uint16(r)
	i := bsearch16(isGraphic, rr)

	return i < len(isGraphic) && rr == isGraphic[i]
}

// bsearch16 returns the smallest i such that a[i] >= x.
// If there is no such i, bsearch16 returns len(a).
func bsearch16(a []uint16, x uint16) int {
	i, j := 0, len(a)

	for i < j {
		h := i + (j-i)>>1
		if a[h] < x {
			i = h + 1
		} else {
			j = h
		}
	}

	return i
}

// isGraphic lists the graphic runes not matched by IsPrint.
var isGraphic = []uint16{
	0x00a0,
	0x1680,
	0x2000,
	0x2001,
	0x2002,
	0x2003,
	0x2004,
	0x2005,
	0x2006,
	0x2007,
	0x2008,
	0x2009,
	0x200a,
	0x202f,
	0x205f,
	0x3000,
}
