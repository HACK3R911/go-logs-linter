package rules

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckSpecialChars(t *testing.T) {
	tests := []struct {
		name           string
		msg            string
		wantDiagnostic bool
	}{
		{
			name:           "normal text",
			msg:            "hello world",
			wantDiagnostic: false,
		},
		{
			name:           "with emoji",
			msg:            "hello 🚀 world",
			wantDiagnostic: true,
		},
		{
			name:           "with fire emoji",
			msg:            "test 🔥 message",
			wantDiagnostic: true,
		},
		{
			name:           "with skull emoji",
			msg:            "error 💀 occurred",
			wantDiagnostic: true,
		},
		{
			name:           "multiple emojis",
			msg:            "🔥🚒💀",
			wantDiagnostic: true,
		},
		{
			name:           "control character",
			msg:            "test\x00message",
			wantDiagnostic: true,
		},
		{
			name:           "tab character",
			msg:            "test\tmessage",
			wantDiagnostic: true,
		},
		{
			name:           "newline character",
			msg:            "test\nmessage",
			wantDiagnostic: true,
		},
		{
			name:           "empty message",
			msg:            "",
			wantDiagnostic: false,
		},
		{
			name:           "punctuation only",
			msg:            "hello, world!",
			wantDiagnostic: false,
		},
		{
			name:           "numbers only",
			msg:            "12345",
			wantDiagnostic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diag := CheckSpecialChars(tt.msg, token.NoPos)
			if tt.wantDiagnostic {
				assert.NotNil(t, diag)
				assert.Equal(t, "style", diag.Category)
			} else {
				assert.Nil(t, diag)
			}
		})
	}
}

func TestCheckSpecialChars_Position(t *testing.T) {
	diag := CheckSpecialChars("test 🚀 emoji", 15)

	assert.NotNil(t, diag)
	assert.Equal(t, token.Pos(15), diag.Pos)
}
