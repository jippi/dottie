package token

import "strconv"

// Type is the set of lexical tokens.
type Type uint

// The list of tokens.
const (
	Illegal Type = iota
	EOF

	//
	// Special characters
	//

	GroupBanner       // # -- ### (3 or more hashtags)
	Comment           // # -- # <anything>
	CommentAnnotation // # -- # @<name> <value>
	Assign            // = -- KEY=VALUE

	//
	// The following tokens are related to variable assignments..
	//

	Identifier // Name of the variable
	Value      // Value is an interpreted value of the variable, if it contains special characters, they will be escaped
	RawValue   // RawValue is used as-is. Special characters are not escaped.
	Space      // All whitespace symbols except \n (new line)
	NewLine    // A new line symbol (\n)
)

var tokens = []string{
	Illegal: "Illegal",
	EOF:     "EOF",

	//
	// Special characters
	//

	GroupBanner:       "GROUP_HEADER",
	Comment:           "COMMENT",
	CommentAnnotation: "COMMENT_ANNOTATION",
	Assign:            "ASSIGN",

	//
	// The following tokens are related to variable assignments..
	//

	Identifier: "IDENTIFIER",
	Value:      "VALUE",
	RawValue:   "RAW_VALUE",
	Space:      "SPACE",
	NewLine:    "NEW_LINE",
}

// String returns the string corresponding to the token.
func (t Type) String() string {
	s := ""

	if int(t) < len(tokens) {
		s = tokens[t]
	}

	if s == "" {
		s = "token(" + strconv.Itoa(int(t)) + ")"
	}

	return s
}
