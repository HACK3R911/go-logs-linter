package rules

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckEnglish(t *testing.T) {
	tests := []struct {
		name           string
		msg            string
		wantDiagnostic bool
	}{
		{
			name:           "all english letters",
			msg:            "hello world",
			wantDiagnostic: false,
		},
		{
			name:           "english with numbers",
			msg:            "user123 logged in",
			wantDiagnostic: false,
		},
		{
			name:           "english with punctuation",
			msg:            "hello, world!",
			wantDiagnostic: false,
		},
		{
			name:           "russian characters",
			msg:            "привет мир",
			wantDiagnostic: true,
		},
		{
			name:           "chinese characters",
			msg:            "你好世界",
			wantDiagnostic: true,
		},
		{
			name:           "japanese characters",
			msg:            "こんにちは",
			wantDiagnostic: true,
		},
		{
			name:           "arabic characters",
			msg:            "مرحبا بالعالم",
			wantDiagnostic: true,
		},
		{
			name:           "mixed english and russian",
			msg:            "hello привет",
			wantDiagnostic: true,
		},
		{
			name:           "empty message",
			msg:            "",
			wantDiagnostic: false,
		},
		{
			name:           "only spaces",
			msg:            "   ",
			wantDiagnostic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diag := CheckEnglish(tt.msg, token.NoPos)
			if tt.wantDiagnostic {
				assert.NotNil(t, diag)
				assert.Equal(t, "language", diag.Category)
			} else {
				assert.Nil(t, diag)
			}
		})
	}
}

func TestCheckEnglish_Position(t *testing.T) {
	diag := CheckEnglish("привет", 25)

	assert.NotNil(t, diag)
	assert.Equal(t, token.Pos(25), diag.Pos)
}
