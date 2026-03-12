package rules

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaseRule_Check(t *testing.T) {
	tests := []struct {
		name           string
		msg            string
		allowUpper     bool
		wantDiagnostic bool
	}{
		{
			name:           "lowercase start - allowed",
			msg:            "user logged in",
			allowUpper:     false,
			wantDiagnostic: false,
		},
		{
			name:           "uppercase start - not allowed",
			msg:            "User logged in",
			allowUpper:     false,
			wantDiagnostic: true,
		},
		{
			name:           "uppercase start - allowed",
			msg:            "User logged in",
			allowUpper:     true,
			wantDiagnostic: false,
		},
		{
			name:           "empty message",
			msg:            "",
			allowUpper:     false,
			wantDiagnostic: false,
		},
		{
			name:           "number start",
			msg:            "123 users",
			allowUpper:     false,
			wantDiagnostic: false,
		},
		{
			name:           "special char start",
			msg:            "_private",
			allowUpper:     false,
			wantDiagnostic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewCaseRule(tt.allowUpper)
			diag := r.Check(tt.msg, token.NoPos)
			if tt.wantDiagnostic {
				assert.NotNil(t, diag)
			} else {
				assert.Nil(t, diag)
			}
		})
	}
}

func TestCaseRule_Check_Position(t *testing.T) {
	r := NewCaseRule(false)
	diag := r.Check("Bad message", 10)

	assert.NotNil(t, diag)
	assert.Equal(t, token.Pos(10), diag.Pos)
}
