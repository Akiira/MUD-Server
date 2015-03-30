// FormattedString
package main

import (
	"github.com/daviddengcn/go-colortext"
)

type FormattedString struct {
	Color ct.Color
	Value string
}

func newFormattedString(msg string) []FormattedString {
	fs := make([]FormattedString, 1, 1)
	fs = append(fs, FormattedString{Color: ct.White, Value: msg})
	return fs
}

func newFormattedString2(color ct.Color, msg string) []FormattedString {
	fs := make([]FormattedString, 1, 1)
	fs = append(fs, FormattedString{Color: color, Value: msg})
	return fs
}

func addMessageToSplice(splice []FormattedString, color ct.Color, msg string) []FormattedString {
	temp := FormattedString{Color: color, Value: msg}

	return append(splice, temp)
}
