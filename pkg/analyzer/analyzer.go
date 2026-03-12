package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"

	"github.com/HACK3R911/go-logs-linter/pkg/analyzer/detector"
	"github.com/HACK3R911/go-logs-linter/pkg/analyzer/rules"
	"github.com/spf13/viper"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type Settings struct {
	AllowUppercaseStart bool     `yaml:"allow_uppercase_start"`
	AllowedPatterns     []string `yaml:"allowed_patterns"`
	DisallowedPatterns  []string `yaml:"disallowed_patterns"`
	AllowNonEnglish     bool     `yaml:"allow_non_english"`
	AllowSpecialChars   bool     `yaml:"allow_special_chars"`
	AllowSensitiveData  bool     `yaml:"allow_sensitive_data"`
	SensitiveKeywords   []string `yaml:"sensitive_keywords"`
}

func DefaultSettings() Settings {
	return Settings{
		AllowUppercaseStart: false,
		AllowedPatterns:     []string{},
		DisallowedPatterns:  []string{},
		AllowNonEnglish:     false,
		AllowSpecialChars:   false,
		AllowSensitiveData:  false,
		SensitiveKeywords:   []string{"password", "token", "secret", "key", "api_key", "apiKey"},
	}
}

func LoadConfig(path string) (*Settings, error) {
	if path == "" {
		return nil, nil
	}
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	var settings Settings
	if v.IsSet("rules.allow_uppercase_start") {
		settings.AllowUppercaseStart = v.GetBool("rules.allow_uppercase_start")
	}
	if v.IsSet("rules.allowed_patterns") {
		settings.AllowedPatterns = v.GetStringSlice("rules.allowed_patterns")
	}
	if v.IsSet("rules.disallowed_patterns") {
		settings.DisallowedPatterns = v.GetStringSlice("rules.disallowed_patterns")
	}
	if v.IsSet("rules.allow_non_english") {
		settings.AllowNonEnglish = v.GetBool("rules.allow_non_english")
	}
	if v.IsSet("rules.allow_special_chars") {
		settings.AllowSpecialChars = v.GetBool("rules.allow_special_chars")
	}
	if v.IsSet("rules.allow_sensitive_data") {
		settings.AllowSensitiveData = v.GetBool("rules.allow_sensitive_data")
	}
	if v.IsSet("rules.sensitive_keywords") {
		settings.SensitiveKeywords = v.GetStringSlice("rules.sensitive_keywords")
	}
	return &settings, nil
}

type logLintAnalyzer struct {
	settings  Settings
	detector  *detector.LoggerDetector
	caseRule  *rules.CaseRule
	inspector *inspector.Inspector
}

func NewAnalyzer(settings Settings) *analysis.Analyzer {
	a := &logLintAnalyzer{
		settings: settings,
	}
	a.detector = detector.NewLoggerDetector()
	a.caseRule = rules.NewCaseRule(settings.AllowUppercaseStart)
	return &analysis.Analyzer{
		Name:     "loglint",
		Doc:      "Checks log messages follow style rules",
		Run:      a.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func (a *logLintAnalyzer) run(pass *analysis.Pass) (any, error) {
	a.inspector = pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	a.inspector.Preorder(nodeFilter, func(node ast.Node) {
		call := node.(*ast.CallExpr)
		a.checkLogCall(pass, call)
	})

	return nil, nil
}

func (a *logLintAnalyzer) checkLogCall(pass *analysis.Pass, call *ast.CallExpr) {
	logCall := a.detector.Detect(pass, call)
	if logCall == nil {
		return
	}

	for _, msg := range logCall.Messages {
		if msg.Text == "" {
			continue
		}

		diagnostic := a.caseRule.Check(msg.Text, msg.Position)
		if diagnostic != nil {
			pass.Report(*diagnostic)
		}

		if !a.settings.AllowNonEnglish {
			diagnostic = rules.CheckEnglish(msg.Text, msg.Position)
			if diagnostic != nil {
				pass.Report(*diagnostic)
			}
		}

		if !a.settings.AllowSpecialChars {
			diagnostic = rules.CheckSpecialChars(msg.Text, msg.Position)
			if diagnostic != nil {
				pass.Report(*diagnostic)
			}
		}

		if !a.settings.AllowSensitiveData {
			diagnostic = rules.CheckSensitiveData(msg.Text, msg.Position, a.settings.SensitiveKeywords)
			if diagnostic != nil {
				pass.Report(*diagnostic)
			}
		}

		for _, pattern := range a.settings.DisallowedPatterns {
			if match, _ := regexp.MatchString(pattern, msg.Text); match {
				pass.Report(analysis.Diagnostic{
					Pos:      msg.Position,
					End:      token.Pos(int(msg.Position) + len(msg.Text)),
					Category: "disallowed",
					Message:  fmt.Sprintf("log message matches disallowed pattern: %s", pattern),
				})
			}
		}

		for _, pattern := range a.settings.AllowedPatterns {
			re := regexp.MustCompile(pattern)
			if !re.MatchString(msg.Text) {
				pass.Report(analysis.Diagnostic{
					Pos:      msg.Position,
					End:      token.Pos(int(msg.Position) + len(msg.Text)),
					Category: "allowed",
					Message:  fmt.Sprintf("log message does not match allowed pattern: %s", pattern),
				})
			}
		}
	}
}
