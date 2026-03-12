package rules

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckSensitiveData(t *testing.T) {
	tests := []struct {
		name           string
		msg            string
		customKeywords []string
		wantDiagnostic bool
	}{
		{
			name:           "normal message",
			msg:            "user logged in successfully",
			customKeywords: nil,
			wantDiagnostic: false,
		},
		{
			name:           "password keyword",
			msg:            "password is incorrect",
			customKeywords: nil,
			wantDiagnostic: true,
		},
		{
			name:           "token keyword",
			msg:            "token expired",
			customKeywords: nil,
			wantDiagnostic: true,
		},
		{
			name:           "secret keyword",
			msg:            "secret key missing",
			customKeywords: nil,
			wantDiagnostic: true,
		},
		{
			name:           "api_key keyword",
			msg:            "api_key not found",
			customKeywords: nil,
			wantDiagnostic: true,
		},
		{
			name:           "email pattern",
			msg:            "contact user@example.com",
			customKeywords: nil,
			wantDiagnostic: true,
		},
		{
			name:           "credit card pattern",
			msg:            "card 4111111111111111",
			customKeywords: nil,
			wantDiagnostic: true,
		},
		{
			name:           "password equals pattern",
			msg:            "password=secret123",
			customKeywords: nil,
			wantDiagnostic: true,
		},
		{
			name:           "bearer token pattern",
			msg:            "bearer abc123token",
			customKeywords: nil,
			wantDiagnostic: true,
		},
		{
			name:           "custom keyword",
			msg:            "found my_secret here",
			customKeywords: []string{"my_secret"},
			wantDiagnostic: true,
		},
		{
			name:           "custom keyword not present",
			msg:            "normal message",
			customKeywords: []string{"my_secret"},
			wantDiagnostic: false,
		},
		{
			name:           "custom keyword not present",
			msg:            "normal message",
			customKeywords: []string{"my_secret"},
			wantDiagnostic: false,
		},
		{
			name:           "empty message",
			msg:            "",
			customKeywords: nil,
			wantDiagnostic: false,
		},
		{
			name:           "case insensitive password",
			msg:            "PASSWORD is secret",
			customKeywords: nil,
			wantDiagnostic: true,
		},
		{
			name:           "auth keyword",
			msg:            "auth failed",
			customKeywords: nil,
			wantDiagnostic: true,
		},
		{
			name:           "credential keyword",
			msg:            "invalid credentials",
			customKeywords: nil,
			wantDiagnostic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diag := CheckSensitiveData(tt.msg, token.NoPos, tt.customKeywords)
			if tt.wantDiagnostic {
				assert.NotNil(t, diag)
				assert.Equal(t, "security", diag.Category)
			} else {
				assert.Nil(t, diag)
			}
		})
	}
}

func TestCheckSensitiveData_Position(t *testing.T) {
	diag := CheckSensitiveData("password is secret", 20, nil)

	assert.NotNil(t, diag)
	assert.Equal(t, token.Pos(20), diag.Pos)
}

func TestDefaultSensitiveKeywords(t *testing.T) {
	assert.NotEmpty(t, DefaultSensitiveKeywords)
	assert.Contains(t, DefaultSensitiveKeywords, "password")
	assert.Contains(t, DefaultSensitiveKeywords, "secret")
	assert.Contains(t, DefaultSensitiveKeywords, "token")
	assert.Contains(t, DefaultSensitiveKeywords, "api_key")
}

func TestSensitivePatterns(t *testing.T) {
	assert.NotEmpty(t, SensitivePatterns)

	testCases := []string{
		"password=secret123",
		"token=abc123",
		"api_key=sk_test",
		"bearer abc123",
		"user@example.com",
		"4111111111111111",
	}

	for _, tc := range testCases {
		found := false
		for _, p := range SensitivePatterns {
			if p.MatchString(tc) {
				found = true
				break
			}
		}
		assert.True(t, found, "pattern should match: %s", tc)
	}
}
