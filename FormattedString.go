// FormattedString
package main

import (
	"github.com/daviddengcn/go-colortext"
)

type FmtStrCollection struct {
	fmtedStrings []FormattedString
	currentColor ct.Color
}

func newFormattedStringCollection() *FmtStrCollection {
	fsc := new(FmtStrCollection)
	fsc.fmtedStrings = make([]FormattedString, 0, 1)
	fsc.currentColor = ct.White
	return fsc
}

func (fsc *FmtStrCollection) addMessage(color ct.Color, msg string) {
	fsc.fmtedStrings = append(fsc.fmtedStrings, newFormattedString2(color, msg))
}

func (fsc *FmtStrCollection) addMessage2(msg string) {
	fsc.fmtedStrings = append(fsc.fmtedStrings, newFormattedString2(ct.White, msg))
}

type FormattedString struct {
	Color ct.Color
	Value string
}

func newFormattedString(msg string) FormattedString {
	return FormattedString{Color: ct.White, Value: msg}
}
func newFormattedString2(color ct.Color, msg string) FormattedString {
	return FormattedString{Color: color, Value: msg}
}

func newFormattedStringSplice(msg string) []FormattedString {
	fs := make([]FormattedString, 0, 1)
	fs = append(fs, FormattedString{Color: ct.White, Value: msg})
	return fs
}

func newFormattedStringSplice2(color ct.Color, msg string) []FormattedString {
	fs := make([]FormattedString, 0, 1)
	fs = append(fs, FormattedString{Color: color, Value: msg})
	return fs
}

func addMessageToSplice(splice []FormattedString, msg string) []FormattedString {
	temp := FormattedString{Color: ct.White, Value: msg}

	return append(splice, temp)
}

func addMessageToSplice2(splice []FormattedString, color ct.Color, msg string) []FormattedString {
	temp := FormattedString{Color: color, Value: msg}

	return append(splice, temp)
}
