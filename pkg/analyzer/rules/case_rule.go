package rules

import (
	"go/token"

	"golang.org/x/tools/go/analysis"
)

type CaseRule struct {
	allowUppercaseStart bool
}

func NewCaseRule(allowUppercaseStart bool) *CaseRule {
	return &CaseRule{
		allowUppercaseStart: allowUppercaseStart,
	}
}

func (r *CaseRule) Check(msg string, pos token.Pos) *analysis.Diagnostic {
	if r.allowUppercaseStart {
		return nil
	}

	if len(msg) == 0 {
		return nil
	}

	firstChar := rune(msg[0])
	if firstChar >= 'A' && firstChar <= 'Z' {
		return &analysis.Diagnostic{
			Pos:      pos,
			End:      token.Pos(int(pos) + 1),
			Category: "style",
			Message:  "log message should start with lowercase letter",
		}
	}

	return nil
}
