package render

type HandlerSignal uint

const (
	Continue HandlerSignal = iota
	Stop
	Return
)

var signals = []string{
	Continue: "CONTINUE",
	Stop:     "STOP",
	Return:   "RETURN",
}

// String returns the string corresponding to the token.
func (hs HandlerSignal) String() string {
	s := ""

	if int(hs) < len(signals) {
		s = signals[hs]
	}

	return s
}
