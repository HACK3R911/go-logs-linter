package main

import (
	"flag"

	"github.com/HACK3R911/go-logs-linter/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var configPath = flag.String("config", "", "path to config.yaml file")

func main() {
	flag.Parse()

	settings := analyzer.DefaultSettings()

	if configPath != nil && *configPath != "" {
		cfg, err := analyzer.LoadConfig(*configPath)
		if err != nil {
		} else if cfg != nil {
			settings = *cfg
		}
	}
	singlechecker.Main(analyzer.NewAnalyzer(settings))
}
