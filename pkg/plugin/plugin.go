package plugin

import (
	"github.com/HACK3R911/go-logs-linter/pkg/analyzer"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("loglint", New)
}

type LogLintPlugin struct {
	settings analyzer.Settings
}

func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[analyzer.Settings](settings)
	if err != nil {
		return nil, err
	}
	return &LogLintPlugin{settings: s}, nil
}

func (p *LogLintPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	analyzer := analyzer.NewAnalyzer(p.settings)
	return []*analysis.Analyzer{analyzer}, nil
}

func (p *LogLintPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
