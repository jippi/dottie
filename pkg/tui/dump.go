package tui

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	slogctx "github.com/veqryn/slog-context"
)

func DumpSlice(ctx context.Context, value string) []string {
	x, _ := Unquote(ctx, value, '"', true) //nolint

	var res []string

	res = append(res, fmt.Sprint("Raw ........  : ", value))
	res = append(res, fmt.Sprint("Glyph ......  : ", fmt.Sprintf("%q", value)))
	res = append(res, fmt.Sprint("UTF-8 ......  : ", fmt.Sprintf("% x", []rune(value))))
	res = append(res, fmt.Sprint("Unicode ....  : ", fmt.Sprintf("%U", []rune(value))))
	res = append(res, fmt.Sprint("TUI Quote ... : ", fmt.Sprintf("%U", []rune(Quote(ctx, value)))))
	res = append(res, fmt.Sprint("TUI Unquote . : ", fmt.Sprintf("%U", []rune(x))))
	res = append(res, fmt.Sprint("[]rune ...... : ", fmt.Sprintf("%v", []rune(value))))
	res = append(res, fmt.Sprint("[]byte ...... : ", fmt.Sprintf("%v", []byte(value))))
	res = append(res, fmt.Sprint("Spew .......  : ", spew.Sdump(value)))

	res = append(res, hex.Dump([]byte(value)))

	return res
}

func Dump(ctx context.Context, value string) {
	for _, line := range DumpSlice(ctx, value) {
		slogctx.Debug(ctx, line)
	}
}
