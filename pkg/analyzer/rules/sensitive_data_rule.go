package rules

import (
	"go/token"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var DefaultSensitiveKeywords = []string{
	"password",
	"passwd",
	"secret",
	"token",
	"api_key",
	"apikey",
	"api-key",
	"access_token",
	"access-token",
	"refresh_token",
	"refresh-token",
	"auth",
	"authorization",
	"private_key",
	"private-key",
	"privatekey",
	"public_key",
	"public-key",
	"publickey",
	"credential",
	"credentials",
	"ssn",
	"social_security",
	"credit_card",
	"creditcard",
	"cvv",
	"cvc",
	"pin",
	"security_code",
	"session_id",
	"sessionid",
	"jwt",
	"bearer",
	"basic_auth",
	"basicauth",
	"client_secret",
	"client-secret",
}

var SensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|passwd|pwd)\s*[:=]\s*\S+`),
	regexp.MustCompile(`(?i)(secret|token)\s*[:=]\s*\S+`),
	regexp.MustCompile(`(?i)(api[_-]?key)\s*[:=]\s*\S+`),
	regexp.MustCompile(`(?i)(bearer|jwt)\s+\S+`),
	regexp.MustCompile(`(?i)(basic\s+auth)\s*[:=]\s*\S+`),
	regexp.MustCompile(`(?i)(private[_-]?key)\s*[:=]\s*\S+`),
	regexp.MustCompile(`(?i)(access[_-]?token)\s*[:=]\s*\S+`),
	regexp.MustCompile(`(?i)(refresh[_-]?token)\s*[:=]\s*\S+`),
	regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`),
	regexp.MustCompile(`\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b`),
}

func CheckSensitiveData(msg string, pos token.Pos, customKeywords []string) *analysis.Diagnostic {
	keywords := DefaultSensitiveKeywords
	if len(customKeywords) > 0 {
		keywords = customKeywords
	}

	lowerMsg := strings.ToLower(msg)

	for _, keyword := range keywords {
		if strings.Contains(lowerMsg, keyword) {
			return &analysis.Diagnostic{
				Pos:      pos,
				End:      token.Pos(int(pos) + len(msg)),
				Category: "security",
				Message:  "log message may contain sensitive data: '" + keyword + "'",
			}
		}
	}

	for _, pattern := range SensitivePatterns {
		if pattern.MatchString(msg) {
			return &analysis.Diagnostic{
				Pos:      pos,
				End:      token.Pos(int(pos) + len(msg)),
				Category: "security",
				Message:  "log message may contain sensitive data pattern",
			}
		}
	}

	return nil
}
