package token

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	slogctx "github.com/veqryn/slog-context"
)

func DebugStringSlice(ctx context.Context, value string) []string {
	x, _ := Unescape(ctx, value, DoubleQuote) //nolint

	var res []string

	res = append(res, fmt.Sprint("Raw ........  : ", value))
	res = append(res, fmt.Sprint("Glyph ......  : ", fmt.Sprintf("%q", value)))
	res = append(res, fmt.Sprint("UTF-8 ......  : ", fmt.Sprintf("% x", []rune(value))))
	res = append(res, fmt.Sprint("Unicode ....  : ", fmt.Sprintf("%U", []rune(value))))
	res = append(res, fmt.Sprint("TUI Quote ... : ", fmt.Sprintf("%U", []rune(Escape(ctx, value, DoubleQuote)))))
	res = append(res, fmt.Sprint("TUI Unquote . : ", fmt.Sprintf("%U", []rune(x))))
	res = append(res, fmt.Sprint("[]rune ...... : ", fmt.Sprintf("%v", []rune(value))))
	res = append(res, fmt.Sprint("[]byte ...... : ", fmt.Sprintf("%v", []byte(value))))
	res = append(res, fmt.Sprint("Spew .......  : ", spew.Sdump(value)))

	res = append(res, hex.Dump([]byte(value)))

	return res
}

func DebugString(ctx context.Context, value string) {
	for _, line := range DebugStringSlice(ctx, value) {
		slogctx.Debug(ctx, line)
	}
}
