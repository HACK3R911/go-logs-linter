package rules

import (
	"go/token"

	"golang.org/x/tools/go/analysis"
)

func CheckSpecialChars(msg string, pos token.Pos) *analysis.Diagnostic {
	for _, r := range msg {
		if isEmoji(r) || isSpecialChar(r) {
			return &analysis.Diagnostic{
				Pos:      pos,
				End:      token.Pos(int(pos) + len(msg)),
				Category: "style",
				Message:  "log message should not contain special characters or emojis",
			}
		}
	}
	return nil
}

func isEmoji(r rune) bool {
	switch {
	case r >= 0x1F600 && r <= 0x1F64F:
		return true
	case r >= 0x1F300 && r <= 0x1F5FF:
		return true
	case r >= 0x1F680 && r <= 0x1F6FF:
		return true
	case r >= 0x1F1E0 && r <= 0x1F1FF:
		return true
	case r >= 0x2700 && r <= 0x27BF:
		return true
	case r >= 0xFE00 && r <= 0xFE0F:
		return true
	case r >= 0x1F900 && r <= 0x1F9FF:
		return true
	case r >= 0x1FA00 && r <= 0x1FA6F:
		return true
	case r >= 0x1FA70 && r <= 0x1FAFF:
		return true
	case r >= 0x2600 && r <= 0x26FF:
		return true
	case r >= 0xFE0F:
		return true
	case r == 0x200D:
		return true
	case r == 0x2640 || r == 0x2642 || r == 0x2695:
		return true
	}
	return false
}

func isSpecialChar(r rune) bool {
	if r >= 0x0000 && r <= 0x001F {
		return true
	}
	if r == 0x007F {
		return true
	}
	return false
}
