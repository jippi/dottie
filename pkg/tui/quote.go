package tui

import (
	"context"
	"strconv"
	"unicode/utf8"

	slogctx "github.com/veqryn/slog-context"
)

func bleh() {
	strconv.Quote("")
}

// Quote quotes each argument and joins them with a space.
// If passed to /bin/sh, the resulting string will be split back into the
// original arguments.
func Quote(ctx context.Context, in string) string {
	return string(quote(ctx, nil, in, '"'))
}

const (
	lowerhex = "0123456789abcdef"
	upperhex = "0123456789ABCDEF"
)

func quote(ctx context.Context, buf []byte, word string, quote byte) []byte {
	slogctx.Debug(ctx, "quote.input.word", StringDump(word))

	var (
		ASCIIonly   = false
		graphicOnly = false
	)

	cur := word

	for width := 0; len(cur) > 0; cur = cur[width:] { //nolint
		runeValue := rune(cur[0])
		width = 1

		if runeValue >= utf8.RuneSelf {
			runeValue, width = utf8.DecodeRuneInString(cur)
		}

		if width == 1 && runeValue == utf8.RuneError {
			slogctx.Debug(ctx, "quote.for-loop.outcome: width == 1 && runeValue == utf8.RuneError")

			buf = append(buf, `\x`...)
			buf = append(buf, lowerhex[cur[0]>>4])
			buf = append(buf, lowerhex[cur[0]&0xF])

			continue
		}

		slogctx.Debug(ctx, "quote.for-loop.outcome: appendEscapedRune")

		buf = appendEscapedRune(ctx, buf, runeValue, quote, ASCIIonly, graphicOnly)
	}

	return buf
}

func appendEscapedRune(ctx context.Context, buf []byte, runeValue rune, quote byte, ASCIIonly, graphicOnly bool) []byte {
	slogctx.Debug(ctx, "appendEscapedRune.input.rune", StringDump(string(runeValue)))

	if runeValue == rune(quote) || runeValue == '\\' { // always backslashed
		slogctx.Debug(ctx, "appendEscapedRune.input.rune: r == rune(quote)")

		buf = append(buf, '\\')
		buf = append(buf, byte(runeValue))

		return buf
	}

	if ASCIIonly {
		slogctx.Debug(ctx, "appendEscapedRune.input.rune: ASCIIonly")

		if runeValue < utf8.RuneSelf && strconv.IsPrint(runeValue) {
			buf = append(buf, byte(runeValue))

			return buf
		}
	} else if strconv.IsPrint(runeValue) || graphicOnly && isInGraphicList(runeValue) {
		slogctx.Debug(ctx, "appendEscapedRune.input.rune: IsPrint/isInGraphicList")

		return utf8.AppendRune(buf, runeValue)
	}

	switch runeValue {
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
		case runeValue < ' ' || runeValue == 0x7f:
			buf = append(buf, `\x`...)
			buf = append(buf, lowerhex[byte(runeValue)>>4])
			buf = append(buf, lowerhex[byte(runeValue)&0xF])
		case !utf8.ValidRune(runeValue):
			runeValue = 0xFFFD

			fallthrough
		case runeValue < 0x10000:
			buf = append(buf, `\u`...)
			for s := 12; s >= 0; s -= 4 {
				buf = append(buf, lowerhex[runeValue>>uint(s)&0xF])
			}
		default:
			buf = append(buf, `\U`...)
			for s := 28; s >= 0; s -= 4 {
				buf = append(buf, lowerhex[runeValue>>uint(s)&0xF])
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
