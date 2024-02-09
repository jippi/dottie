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

// String returns the string corresponding to the Handler Signal.
func (hs HandlerSignal) String() string {
	str := ""

	if int(hs) < len(signals) {
		str = signals[hs]
	}

	return str
}
