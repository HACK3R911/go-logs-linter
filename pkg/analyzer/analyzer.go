package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"

	"github.com/HACK3R911/go-logs-linter/pkg/analyzer/detector"
	"github.com/HACK3R911/go-logs-linter/pkg/analyzer/rules"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type Settings struct {
	AllowUppercaseStart bool
	AllowedPatterns     []string
	DisallowedPatterns  []string
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
