package set_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/davecgh/go-spew/spew"
	"github.com/jippi/dottie/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func FuzzSetCommand(f *testing.F) {
	f.Add("@@\v\"@23")

	f.Fuzz(doTest)
}

func TestSpecificInputs(t *testing.T) {
	t.Parallel()

	t.Run("newline", func(t *testing.T) {
		t.Parallel()

		doTest(t, "\n")
	})

	t.Run("tab", func(t *testing.T) {
		t.Parallel()

		doTest(t, "\t")
	})

	t.Run("slash", func(t *testing.T) {
		t.Parallel()

		doTest(t, "\\")
	})

	t.Run("null", func(t *testing.T) {
		t.Parallel()

		doTest(t, "\x00")
	})

	t.Run("weird", func(t *testing.T) {
		t.Parallel()

		doTest(t, "@@\v\"@23")
	})

	t.Run("weird_2", func(t *testing.T) {
		t.Parallel()

		doTest(t, "\"\n")
	})
}

func doTest(t *testing.T, expected string) { //nolint thelper
	dotEnvFile := t.TempDir() + "/tmp.env"

	_, err := os.Create(dotEnvFile)
	require.NoErrorf(t, err, "failed to create empty .env file [ %s ] in TempDir", dotEnvFile)

	t.Log("-----------------------")
	t.Log("EXPECTED VALUE")
	t.Log("-----------------------")
	dump(t, expected)

	// Set the KEY/VALUE pair
	setFailed := false

	{
		var (
			stdout bytes.Buffer
			stderr bytes.Buffer
			args   = []string{
				"--file", dotEnvFile,
				"set",
				"--",
				"my_key", expected,
			}
		)

		t.Log("-----------------------")
		t.Log("ARGS:")
		t.Log("-----------------------")
		t.Log(spew.Sdump(args))
		t.Log()

		// Run command
		_, err := cmd.RunCommand(context.Background(), args, &stdout, &stderr)

		if stdout.Len() == 0 {
			stdout.WriteString("(empty)")
		}

		if stderr.Len() == 0 {
			stderr.WriteString("(empty)")
		}

		t.Log("-----------------------")
		t.Log("STDOUT:")
		t.Log("-----------------------")
		t.Log(stdout.String())
		t.Log()

		t.Log("-----------------------")
		t.Log("STDERR:")
		t.Log("-----------------------")
		t.Log(stderr.String())
		t.Log()

		switch err {
		case nil:
			out := stdout.String()

			assert.Contains(t, out, "Key [ my_key ] was successfully upserted")
			assert.Contains(t, out, "File was successfully saved")

		default:
			setFailed = true

			assert.Regexp(t, "(illegal UTF-8 encoding|Invalid template)", stderr.String())
		}
	}

	// Half-way checkpoint for some checks

	{
		if setFailed {
			return
		}

		out, err := os.ReadFile(dotEnvFile)
		require.NoError(t, err)

		disk := strings.TrimRight(string(out), "\n")

		t.Log("-----------------------")
		t.Log("FILE ON DISK")
		t.Log("-----------------------")
		dump(t, disk)
	}

	// Read back from disk
	{
		var (
			stdout      bytes.Buffer
			stderr      bytes.Buffer
			commandArgs = []string{
				"--file", dotEnvFile,
				"value",
				// "--literal",
				"my_key",
			}
		)

		// Run command
		cmd.RunCommand(context.Background(), commandArgs, &stdout, &stderr)

		t.Log("-----------------------")
		t.Log("STDOUT:")
		t.Log("-----------------------")
		t.Log(stdout.String())

		t.Log("-----------------------")
		t.Log("STDERR:")
		t.Log("-----------------------")
		t.Log(stderr.String())

		// In cases where we get "$0" and similar, the actual interpolated output will be different than the input
		if strings.Contains(stderr.String(), "Defaulting to a blank string") {
			return
		}

		actual := stdout.String()
		// actual = token.DoubleQuotes.Escape(actual)
		// actual = Clean(actual)

		t.Log("-----------------------")
		t.Log("Actual")
		t.Log("-----------------------")
		dump(t, actual)

		assert.Equal(t, fmt.Sprintf("%U", []rune(expected)), fmt.Sprintf("%U", []rune(actual)))
		// assert.True(t, false, "fail")
	}
}

func Clean(str string) string {
	return strings.Map(func(runeVal rune) rune {
		if !utf8.ValidRune(runeVal) {
			return -1
		}

		if !unicode.IsGraphic(runeVal) {
			return -1
		}

		return runeVal
	}, str)
}

func dump(t *testing.T, value string) {
	t.Helper()

	t.Log("Raw .....  :", value)
	t.Log("Glyph ...  :", fmt.Sprintf("%q", value))
	t.Log("UTF-8 ...  :", fmt.Sprintf("% x", []rune(value)))
	t.Log("Unicode .  :", fmt.Sprintf("%U", []rune(value)))
	t.Log("Clean A .. :", Clean(value))
	t.Log("Clean B .. :", fmt.Sprintf("%U", []rune(Clean(value))))
	t.Log("quote A..  :", appendQuotedWith(nil, value, false, false))
	t.Log("quote B..  :", appendQuotedWith(nil, value, false, true))
	t.Log("Spew ....  :", spew.Sdump(value))
}

const (
	lowerhex = "0123456789abcdef"
	upperhex = "0123456789ABCDEF"
)

func appendQuotedWith(buf []byte, str string, ASCIIonly, graphicOnly bool) string {
	// Often called with big strings, so preallocate. If there's quoting,
	// this is conservative but still helps a lot.
	if cap(buf)-len(buf) < len(str) {
		nBuf := make([]byte, len(buf), len(buf)+1+len(str)+1)
		copy(nBuf, buf)
		buf = nBuf
	}

	for width := 0; len(str) > 0; str = str[width:] { //nolint
		runeC := rune(str[0])

		width = 1
		if runeC >= utf8.RuneSelf {
			runeC, width = utf8.DecodeRuneInString(str)
		}

		if width == 1 && runeC == utf8.RuneError {
			buf = append(buf, `\x`...)
			buf = append(buf, lowerhex[str[0]>>4])
			buf = append(buf, lowerhex[str[0]&0xF])

			continue
		}

		buf = appendEscapedRune(buf, runeC, ASCIIonly, graphicOnly)
	}

	return string(buf)
}

func appendEscapedRune(buf []byte, runeVal rune, ASCIIonly, graphicOnly bool) []byte {
	if runeVal == '\\' { // always backslashed
		buf = append(buf, '\\')
		buf = append(buf, byte(runeVal))

		return buf
	}

	if ASCIIonly {
		if runeVal < utf8.RuneSelf && strconv.IsPrint(runeVal) {
			buf = append(buf, byte(runeVal))

			return buf
		}
	} else if strconv.IsPrint(runeVal) || graphicOnly && isInGraphicList(runeVal) {
		return utf8.AppendRune(buf, runeVal)
	}

	switch runeVal {
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
		case runeVal < ' ' || runeVal == 0x7f:
			buf = append(buf, `\x`...)
			buf = append(buf, lowerhex[byte(runeVal)>>4])
			buf = append(buf, lowerhex[byte(runeVal)&0xF])
		case !utf8.ValidRune(runeVal):
			runeVal = 0xFFFD

			fallthrough

		case runeVal < 0x10000:
			buf = append(buf, `\u`...)
			for s := 12; s >= 0; s -= 4 {
				buf = append(buf, lowerhex[runeVal>>uint(s)&0xF])
			}

		default:
			buf = append(buf, `\U`...)
			for s := 28; s >= 0; s -= 4 {
				buf = append(buf, lowerhex[runeVal>>uint(s)&0xF])
			}
		}
	}

	return buf
}

// isInGraphicList reports whether the rune is in the isGraphic list. This separation
// from IsGraphic allows quoteWith to avoid two calls to IsPrint.
// Should be called only if IsPrint fails.
func isInGraphicList(runeVal rune) bool {
	// We know r must fit in 16 bits - see makeisprint.go.
	if runeVal > 0xFFFF {
		return false
	}

	rr := uint16(runeVal)

	i := bsearch16(isGraphic, rr)

	return i < len(isGraphic) && rr == isGraphic[i]
}

// bsearch16 returns the smallest i such that a[i] >= x.
// If there is no such i, bsearch16 returns len(a).
func bsearch16(a []uint16, x uint16) int {
	index, length := 0, len(a)

	for index < length {
		h := index + (length-index)>>1
		if a[h] < x {
			index = h + 1
		} else {
			length = h
		}
	}

	return index
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
