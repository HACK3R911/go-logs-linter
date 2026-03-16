package rules

import (
	"go/token"
	"testing"

	"github.com/HACK3R911/go-logs-linter/pkg/analyzer/detector"
	"github.com/stretchr/testify/assert"
)

func TestCheckSensitiveData_Concatenation(t *testing.T) {
	tests := []struct {
		name           string
		logCall        *detector.LogCall
		wantDiagnostic bool
	}{
		{
			name: "password variable in concatenation",
			logCall: &detector.LogCall{
				Messages: []detector.LogMessage{
					{Text: "test", Position: token.NoPos},
				},
				Concatenation: detector.ConcatenationInfo{
					IsConcatenation: true,
					VarNames:        []string{"password"},
				},
			},
			wantDiagnostic: true,
		},
		{
			name: "token variable in concatenation",
			logCall: &detector.LogCall{
				Messages: []detector.LogMessage{
					{Text: "test", Position: token.NoPos},
				},
				Concatenation: detector.ConcatenationInfo{
					IsConcatenation: true,
					VarNames:        []string{"token"},
				},
			},
			wantDiagnostic: true,
		},
		{
			name: "apiKey variable in concatenation",
			logCall: &detector.LogCall{
				Messages: []detector.LogMessage{
					{Text: "test", Position: token.NoPos},
				},
				Concatenation: detector.ConcatenationInfo{
					IsConcatenation: true,
					VarNames:        []string{"apiKey"},
				},
			},
			wantDiagnostic: true,
		},
		{
			name: "secret variable in concatenation",
			logCall: &detector.LogCall{
				Messages: []detector.LogMessage{
					{Text: "test", Position: token.NoPos},
				},
				Concatenation: detector.ConcatenationInfo{
					IsConcatenation: true,
					VarNames:        []string{"secret"},
				},
			},
			wantDiagnostic: true,
		},
		{
			name: "normal variable in concatenation",
			logCall: &detector.LogCall{
				Messages: []detector.LogMessage{
					{Text: "test", Position: token.NoPos},
				},
				Concatenation: detector.ConcatenationInfo{
					IsConcatenation: true,
					VarNames:        []string{"username"},
				},
			},
			wantDiagnostic: false,
		},
		{
			name: "no concatenation",
			logCall: &detector.LogCall{
				Messages: []detector.LogMessage{
					{Text: "test", Position: token.NoPos},
				},
				Concatenation: detector.ConcatenationInfo{
					IsConcatenation: false,
					VarNames:        []string{},
				},
			},
			wantDiagnostic: false,
		},
		{
			name:           "nil logCall",
			logCall:        nil,
			wantDiagnostic: false,
		},
		{
			name: "custom keyword match",
			logCall: &detector.LogCall{
				Messages: []detector.LogMessage{
					{Text: "test", Position: token.NoPos},
				},
				Concatenation: detector.ConcatenationInfo{
					IsConcatenation: true,
					VarNames:        []string{"mySecret"},
				},
			},
			wantDiagnostic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diag := CheckSensitiveData(tt.logCall, nil)
			if tt.wantDiagnostic {
				assert.NotNil(t, diag)
				assert.Equal(t, "security", diag.Category)
			} else {
				assert.Nil(t, diag)
			}
		})
	}
}

func TestCheckSensitiveData_ZapKeys(t *testing.T) {
	tests := []struct {
		name           string
		logCall        *detector.LogCall
		wantDiagnostic bool
	}{
		{
			name: "password zap key",
			logCall: &detector.LogCall{
				Messages: []detector.LogMessage{
					{Text: "test", Position: token.NoPos},
				},
				Concatenation: detector.ConcatenationInfo{
					IsConcatenation: false,
					ZapKeys:         []string{"password"},
				},
			},
			wantDiagnostic: true,
		},
		{
			name: "token zap key",
			logCall: &detector.LogCall{
				Messages: []detector.LogMessage{
					{Text: "test", Position: token.NoPos},
				},
				Concatenation: detector.ConcatenationInfo{
					IsConcatenation: false,
					ZapKeys:         []string{"token"},
				},
			},
			wantDiagnostic: true,
		},
		{
			name: "api_key zap key",
			logCall: &detector.LogCall{
				Messages: []detector.LogMessage{
					{Text: "test", Position: token.NoPos},
				},
				Concatenation: detector.ConcatenationInfo{
					IsConcatenation: false,
					ZapKeys:         []string{"api_key"},
				},
			},
			wantDiagnostic: true,
		},
		{
			name: "secret zap key",
			logCall: &detector.LogCall{
				Messages: []detector.LogMessage{
					{Text: "test", Position: token.NoPos},
				},
				Concatenation: detector.ConcatenationInfo{
					IsConcatenation: false,
					ZapKeys:         []string{"secret"},
				},
			},
			wantDiagnostic: true,
		},
		{
			name: "normal zap key",
			logCall: &detector.LogCall{
				Messages: []detector.LogMessage{
					{Text: "test", Position: token.NoPos},
				},
				Concatenation: detector.ConcatenationInfo{
					IsConcatenation: false,
					ZapKeys:         []string{"username"},
				},
			},
			wantDiagnostic: false,
		},
		{
			name: "multiple zap keys - one sensitive",
			logCall: &detector.LogCall{
				Messages: []detector.LogMessage{
					{Text: "test", Position: token.NoPos},
				},
				Concatenation: detector.ConcatenationInfo{
					IsConcatenation: false,
					ZapKeys:         []string{"username", "password", "email"},
				},
			},
			wantDiagnostic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diag := CheckSensitiveData(tt.logCall, nil)
			if tt.wantDiagnostic {
				assert.NotNil(t, diag)
				assert.Equal(t, "security", diag.Category)
			} else {
				assert.Nil(t, diag)
			}
		})
	}
}

func TestCheckSensitiveData_CustomKeywords(t *testing.T) {
	logCall := &detector.LogCall{
		Messages: []detector.LogMessage{
			{Text: "test", Position: token.NoPos},
		},
		Concatenation: detector.ConcatenationInfo{
			IsConcatenation: true,
			VarNames:        []string{"myCustomSecret"},
		},
	}

	diag := CheckSensitiveData(logCall, []string{"myCustomSecret"})
	assert.NotNil(t, diag)
	assert.Equal(t, "security", diag.Category)

	logCall2 := &detector.LogCall{
		Messages: []detector.LogMessage{
			{Text: "test", Position: token.NoPos},
		},
		Concatenation: detector.ConcatenationInfo{
			IsConcatenation: true,
			VarNames:        []string{"userName"},
		},
	}
	diag = CheckSensitiveData(logCall2, []string{})
	assert.Nil(t, diag)
}

func TestDefaultSensitiveKeywords(t *testing.T) {
	assert.NotEmpty(t, DefaultSensitiveKeywords)
	assert.Contains(t, DefaultSensitiveKeywords, "password")
	assert.Contains(t, DefaultSensitiveKeywords, "token")
	assert.Contains(t, DefaultSensitiveKeywords, "secret")
	assert.Contains(t, DefaultSensitiveKeywords, "api_key")
	assert.Contains(t, DefaultSensitiveKeywords, "jwt")
}
