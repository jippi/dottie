package tui

import (
	"encoding/hex"
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

func DumpSlice(value string) []string {
	x, _ := Unquote(value, '"', true) //nolint

	var res []string

	res = append(res, fmt.Sprint("Raw ........  : ", value))
	res = append(res, fmt.Sprint("Glyph ......  : ", fmt.Sprintf("%q", value)))
	res = append(res, fmt.Sprint("UTF-8 ......  : ", fmt.Sprintf("% x", []rune(value))))
	res = append(res, fmt.Sprint("Unicode ....  : ", fmt.Sprintf("%U", []rune(value))))
	res = append(res, fmt.Sprint("TUI Quote ... : ", fmt.Sprintf("%U", []rune(Quote(value)))))
	res = append(res, fmt.Sprint("TUI Unquote . : ", fmt.Sprintf("%U", []rune(x))))
	res = append(res, fmt.Sprint("[]rune ...... : ", fmt.Sprintf("%v", []rune(value))))
	res = append(res, fmt.Sprint("[]byte ...... : ", fmt.Sprintf("%v", []byte(value))))
	res = append(res, fmt.Sprint("Spew .......  : ", spew.Sdump(value)))

	res = append(res, hex.Dump([]byte(value)))

	return res
}

func Dump(value string) {
	for _, line := range DumpSlice(value) {
		fmt.Println(line)
	}
}
