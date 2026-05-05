package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/pipeline"
	"github.com/yourorg/logslice/internal/source"
)

// Config holds parsed CLI arguments.
type Config struct {
	Filter string
	Format string
	Fields []string
	Files  []string
}

// Run parses args and executes the pipeline.
func Run(args []string) error {
	cfg, err := parseArgs(args)
	if err != nil {
		return err
	}

	f, err := filter.Parse(cfg.Filter)
	if err != nil {
		return fmt.Errorf("invalid filter: %w", err)
	}

	w, err := output.NewWriter(os.Stdout, cfg.Format, cfg.Fields)
	if err != nil {
		return fmt.Errorf("invalid output config: %w", err)
	}

	var src source.Source
	if len(cfg.Files) == 0 {
		src = source.NewStdinSource()
	} else {
		src, err = source.NewFileSource(cfg.Files...)
		if err != nil {
			return fmt.Errorf("opening files: %w", err)
		}
	}

	p := pipeline.New(src, f, w)
	return p.Run()
}

func parseArgs(args []string) (*Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)

	filterExpr := fs.String("filter", "", "filter expression, e.g. 'level=error'")
	format := fs.String("format", "json", fmt.Sprintf("output format: %s", strings.Join(output.ValidFormats, ", ")))
	fields := fs.String("fields", "", "comma-separated list of fields to include in output")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: logslice [options] [file...]\n\nOptions:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	var fieldList []string
	if *fields != "" {
		for _, f := range strings.Split(*fields, ",") {
			if f = strings.TrimSpace(f); f != "" {
				fieldList = append(fieldList, f)
			}
		}
	}

	return &Config{
		Filter: *filterExpr,
		Format: *format,
		Fields: fieldList,
		Files:  fs.Args(),
	}, nil
}
