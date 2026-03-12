package rules

import (
	"go/token"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

func CheckEnglish(msg string, pos token.Pos) *analysis.Diagnostic {
	for _, r := range msg {
		if !isEnglishRune(r) {
			return &analysis.Diagnostic{
				Pos:      pos,
				End:      token.Pos(int(pos) + len(msg)),
				Category: "language",
				Message:  "log message should contain only English (Latin) characters",
			}
		}
	}
	return nil
}

func isEnglishRune(r rune) bool {
	if r <= 0x007F {
		return true
	}
	if r >= 0x00C0 && r <= 0x024F {
		return true
	}
	if unicode.IsSpace(r) {
		return true
	}
	if unicode.IsPunct(r) {
		return true
	}
	if unicode.IsNumber(r) {
		return true
	}
	return false
}
