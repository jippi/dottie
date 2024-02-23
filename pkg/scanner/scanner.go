package scanner

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/jippi/dottie/pkg/token"
)

var escaper = strings.NewReplacer(
	`\n`, "\n",
	`\t`, "\t",
	`\r`, "\r",
	`\v`, "\v",
	`\f`, "\f",
)

const (
	bom = 0xFEFF // byte order mark, only permitted as the first character
	eof = -1     // eof indicates the end of the file.
)

// Scanner converts a sequence of characters into a sequence of tokens.
type Scanner struct {
	input      string
	rune       rune // current character
	prevOffset int  // position before current character
	offset     int  // character offset
	peekOffset int  // position after current character
	lineNumber uint // current line number
}

// New returns new Scanner.
func New(input string) *Scanner {
	scanner := &Scanner{
		input:      input,
		lineNumber: 1,
	}

	scanner.next()

	if scanner.rune == bom {
		scanner.next() // ignore BOM at the beginning of the file
	}

	return scanner
}

// NextToken scans the next token and returns the token position, the token, and its literal string
// if applicable. The source end is indicated by token.EOF.
//
// If the returned token is a literal (token.Identifier, token.Value, token.RawValue) or token.Comment,
// the literal string has the corresponding value.
//
// If the returned token is token.Illegal, the literal string is the offending character.
func (s *Scanner) NextToken() token.Token {
	fmt.Println("Scanner working on rune:", s.rune, string(s.rune))

	switch s.rune {
	case eof:
		return token.New(
			token.EOF,
			token.WithOffset(s.offset),
			token.WithLineNumber(s.lineNumber),
		)

	case '\n':
		return s.scanNewLine()

	case ' ', '\t', '\r', '\v', '\f':
		defer s.next()

		return token.New(
			token.Space,
			token.WithLiteralRune(s.rune),
			token.WithOffset(s.offset),
			token.WithLineNumber(s.lineNumber),
		)

	case '=':
		defer s.next()

		return token.New(
			token.Assign,
			token.WithOffset(s.offset),
			token.WithLineNumber(s.lineNumber),
		)

	case '#':
		return s.scanComment()

	case '"':
		return s.scanQuotedValue(token.Value, token.DoubleQuotes)

	case '\'':
		return s.scanQuotedValue(token.RawValue, token.SingleQuotes)

	default:
		switch prev := s.prev(); prev {
		case '\n', bom:
			if isValidIdentifier(s.rune) {
				return s.scanIdentifier()
			}

		case '=':
			return s.scanUnquotedValue()
		}

		return s.scanIllegalRune()
	}
}

// ========================================================================
// Methods that scan a specific token kind.
// ========================================================================

func (s *Scanner) scanNewLine() token.Token {
	s.lineNumber++

	s.next()

	return token.New(
		token.NewLine,
		token.WithLiteral("\n"),
		token.WithOffset(s.offset),
		token.WithLineNumber(s.lineNumber-1),
	)
}

func (s *Scanner) scanIdentifier() token.Token {
	start := s.offset

	for isLetter(s.rune) || isDigit(s.rune) || isSymbol(s.rune) {
		s.next()
	}

	literal := s.input[start:s.offset]

	return token.New(
		token.Identifier,
		token.WithLiteral(literal),
		token.WithOffset(s.offset),
		token.WithLineNumber(s.lineNumber),
	)
}

func (s *Scanner) scanComment() token.Token {
	start := s.offset

	s.next()

	// If a comment looks like "#KEY=VALUE" it's a commented/disabled KEY=VALUE pair
	// so consume it as such instead of a comment
	if isValidIdentifier(s.rune) {
		res := s.scanIdentifier()
		res.Offset = res.Offset - 1
		res.Commented = true

		return res
	}

	s.skipWhitespace()

	switch {
	// We got an annotation!
	case s.rune == '@':
		return s.scanCommentAnnotation(start)

	// We got a group header
	case s.peek(2) == "##":
		s.untilEndOfLine()

		return token.New(
			token.GroupBanner,
			token.WithLiteral(s.input[start:s.offset]),
			token.WithOffset(s.offset),
			token.WithLineNumber(s.lineNumber),
		)
	}

	s.untilEndOfLine()
	lit := s.input[start:s.offset]

	return token.New(
		token.Comment,
		token.WithLiteral(lit),
		token.WithOffset(s.offset),
		token.WithLineNumber(s.lineNumber),
	)
}

func (s *Scanner) scanCommentAnnotation(offset int) token.Token {
	// Consume the @
	s.next()

	start := s.offset

	// Key
	for !isWideSpace(s.rune) {
		s.next()
	}

	key := s.input[start:s.offset]

	// Consume any space between key and value
	s.skipWhitespace()

	// Value
	valueStart := s.offset
	s.untilEndOfLine()

	value := s.input[valueStart:s.offset]

	return token.New(
		token.CommentAnnotation,
		token.WithLiteral(s.input[offset:s.offset]), // full line
		token.WithOffset(s.offset),
		token.WithLineNumber(s.lineNumber),
		token.WithAnnotation(key, value),
	)
}

func (s *Scanner) skipWhitespace() {
	for !isEOF(s.rune) && !isNewLine(s.rune) && unicode.IsSpace(s.rune) {
		s.next()
	}
}

func (s *Scanner) untilEndOfLine() {
	for !isEOF(s.rune) && !isNewLine(s.rune) {
		s.next()
	}
}

func (s *Scanner) scanIllegalRune() token.Token {
	defer s.next()

	return token.New(
		token.Illegal,
		token.WithLiteralRune(s.rune),
		token.WithOffset(s.offset),
		token.WithLineNumber(s.lineNumber),
	)
}

func (s *Scanner) scanUnquotedValue() token.Token {
	start := s.offset

	for !isEOF(s.rune) && !isNewLine(s.rune) {
		s.next()
	}

	lit := escape(s.input[start:s.offset])

	return token.New(
		token.Value,
		token.WithLiteral(lit),
		token.WithQuoteType(token.NoQuotes),
		token.WithOffset(s.offset),
		token.WithLineNumber(s.lineNumber),
	)
}

func (s *Scanner) scanQuotedValue(tType token.Type, quote token.Quote) token.Token {
	fmt.Println("scanQuotedValue!")
	// opening quote already consumed
	s.next()

	start := s.offset

	fmt.Println("scanQuotedValue: s.input", fmt.Sprintf(">%q<", s.input))

	escapes := 0

	for {
		escapingPrevious := escapes == 1

		fmt.Println("scanQuotedValue -->", fmt.Sprintf("%q", s.rune), fmt.Sprintf("%q", quote.Rune()), s.rune, "inEscape?", escapingPrevious, escapes)

		if isEOF(s.rune) || isNewLine(s.rune) {
			// panic("nein")
			tType = token.Illegal

			break
		}

		// Break parsing if we hit our quote style,
		// and the previous token IS NOT an escape sequence
		if quote.Is(s.rune) && !escapingPrevious {
			// panic("oh no")
			break
		}

		if s.rune == '\\' {
			escapes++
		} else {
			escapes = 0
		}

		if escapes == 2 {
			escapes = 0
		}

		s.next()
	}

	offset := s.offset
	lit := s.input[start:offset]

	fmt.Println("scanQuotedValue lit (before):", tType, lit)

	if tType == token.Value {
		lit = escape(lit)
	}

	fmt.Println("scanQuotedValue lit (after):", tType, lit)

	if quote.Is(s.rune) {
		s.next()
	}

	return token.New(
		tType,
		token.WithLiteral(lit),
		token.WithQuoteType(quote),
		token.WithOffset(offset),
		token.WithLineNumber(s.lineNumber),
	)
}

// ========================================================================
// Methods that control pointers to the current, previous, and next chars.
// ========================================================================

// Read the next Unicode char into s.ch.
// s.ch < 0 means end-of-file.
func (s *Scanner) next() {
	s.prevOffset = s.offset

	if s.peekOffset < len(s.input) {
		s.offset = s.peekOffset
		r, width := s.scanRune(s.offset)

		s.peekOffset += width
		s.rune = r
	} else {
		s.offset = len(s.input)
		s.rune = eof
	}

	if s.offset == 0 {
		s.prevOffset = -1
	}
}

func (s *Scanner) prev() rune {
	switch {
	case s.prevOffset < 0:
		return '\n'

	case s.prevOffset < len(s.input):
		r, _ := s.scanRune(s.prevOffset)

		return r

	default:
		return eof
	}
}

// Reads a single Unicode character and returns the rune and its width in bytes.
func (s *Scanner) scanRune(offset int) (rune, int) {
	runeVal := rune(s.input[offset])
	width := 1

	switch {
	case runeVal >= utf8.RuneSelf:
		// not ASCII
		runeVal, width = utf8.DecodeRuneInString(s.input[offset:])
		if runeVal == utf8.RuneError && width == 1 {
			panic("illegal UTF-8 encoding on position " + strconv.Itoa(offset))
		}

		if runeVal == bom && s.offset > 0 {
			panic("illegal byte order mark on position " + strconv.Itoa(offset))
		}
	}

	return runeVal, width
}

func (s *Scanner) peek(length int) string {
	start := s.offset
	end := start + length

	maxLength := len(s.input)
	if maxLength >= end {
		return s.input[start:end]
	}

	return s.input[start:maxLength]
}

// ========================================================================
// Auxiliary methods that check if the rune is one of the specific kind.
// ========================================================================

func isValidIdentifier(r rune) bool {
	return isLetter(r) || isDigit(r) || isSymbol(r)
}

func isLetter(r rune) bool {
	return ('a' <= lower(r) && lower(r) <= 'z') ||
		(r >= utf8.RuneSelf && unicode.IsLetter(r))
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isSymbol(r rune) bool {
	switch r {
	case '_', '.', ',', '-':
		return true
	}

	return false
}

func isNewLine(r rune) bool {
	return r == '\n'
}

func isEOF(r rune) bool {
	return r == eof
}

func isWideSpace(r rune) bool {
	return unicode.IsSpace(r) || isEOF(r) || isNewLine(r) || r == '\v'
}

// ------------------------------------------------------------------------

// returns lower-case r if r is an ASCII letter
func lower(r rune) rune {
	return ('a' - 'A') | r
}

func escape(s string) string {
	return escaper.Replace(s)
}
