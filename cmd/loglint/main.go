package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/HACK3R911/go-logs-linter/pkg/analyzer"
	"github.com/HACK3R911/go-logs-linter/pkg/config"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var configPath = flag.String("config", "", "path to config.yaml file")

func main() {
	flag.Parse()

	settings := analyzer.DefaultSettings()

	if configPath != nil && *configPath != "" {
		cfg, err := config.Load(*configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		} else if cfg != nil {
			settings = cfg.ToSettings()
		}
	}
	singlechecker.Main(analyzer.NewAnalyzer(settings))
}
