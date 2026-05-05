package pipeline

import (
	"encoding/json"
	"io"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
)

// Pipeline reads JSON log entries from a source, applies a filter,
// and writes matching entries to an output writer.
type Pipeline struct {
	src    io.Reader
	filter *filter.Filter
	writer *output.Writer
}

// New creates a new Pipeline.
func New(src io.Reader, f *filter.Filter, w *output.Writer) *Pipeline {
	return &Pipeline{
		src:    src,
		filter: f,
		writer: w,
	}
}

// Run processes log entries from the source until EOF or an error occurs.
// It returns the number of entries written and any non-EOF error.
func (p *Pipeline) Run() (int, error) {
	dec := json.NewDecoder(p.src)
	written := 0

	for {
		var entry map[string]interface{}
		if err := dec.Decode(&entry); err != nil {
			if err == io.EOF {
				break
			}
			return written, err
		}

		if p.filter != nil && !p.filter.Match(entry) {
			continue
		}

		if err := p.writer.Write(entry); err != nil {
			return written, err
		}
		written++
	}

	return written, nil
}
