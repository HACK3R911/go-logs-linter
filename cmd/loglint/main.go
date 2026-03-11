package main

import (
	"github.com/HACK3R911/go-logs-linter/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.NewAnalyzer(analyzer.Settings{}))
}
