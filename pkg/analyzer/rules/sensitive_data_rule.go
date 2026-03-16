package rules

import (
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/HACK3R911/go-logs-linter/pkg/analyzer/detector"
)

var DefaultSensitiveKeywords = []string{
	"password",
	"passwd",
	"pwd",
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
	"secret_key",
	"secret-key",
	"encryption_key",
	"encryption-key",
}

func CheckSensitiveData(logCall *detector.LogCall, customKeywords []string) *analysis.Diagnostic {
	if logCall == nil {
		return nil
	}

	if len(logCall.Messages) == 0 {
		return nil
	}

	keywords := DefaultSensitiveKeywords
	if len(customKeywords) > 0 {
		keywords = customKeywords
	}

	position := logCall.Messages[0].Position
	concat := logCall.Concatenation

	if concat.IsConcatenation {
		for _, varName := range concat.VarNames {
			lowerVarName := strings.ToLower(varName)
			for _, keyword := range keywords {
				lowerKeyword := strings.ToLower(keyword)
				if strings.Contains(lowerVarName, lowerKeyword) {
					return &analysis.Diagnostic{
						Pos:      position,
						End:      position,
						Category: "security",
						Message:  "log message may contain sensitive data in concatenation: variable '" + varName + "'",
					}
				}
			}
		}
	}

	for _, zapKey := range concat.ZapKeys {
		lowerKey := strings.ToLower(zapKey)
		for _, keyword := range keywords {
			lowerKeyword := strings.ToLower(keyword)
			if strings.Contains(lowerKey, lowerKeyword) {
				return &analysis.Diagnostic{
					Pos:      position,
					End:      position,
					Category: "security",
					Message:  "log message may contain sensitive data in zap key: '" + zapKey + "'",
				}
			}
		}
	}

	return nil
}
