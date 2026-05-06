package tail_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/yourusername/logslice/internal/tail"
)

func writeTmp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestFollow_ReadsExistingLines(t *testing.T) {
	path := writeTmp(t, "{\"level\":\"info\"}\n{\"level\":\"warn\"}\n")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	lines, err := tail.Follow(ctx, path, tail.Options{PollInterval: 20 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}

	var got []string
	for line := range lines {
		got = append(got, line)
		if len(got) == 2 {
			cancel()
		}
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

func TestFollow_PicksUpNewWrites(t *testing.T) {
	path := writeTmp(t, "")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	lines, err := tail.Follow(ctx, path, tail.Options{PollInterval: 20 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
		f, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
		defer f.Close()
		f.WriteString("{\"msg\":\"hello\"}\n")
	}()

	select {
	case line := <-lines:
		if line == "" {
			t.Fatal("expected non-empty line")
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for new line")
	}
}

func TestFollow_InvalidPath(t *testing.T) {
	ctx := context.Background()
	_, err := tail.Follow(ctx, "/nonexistent/path/file.log", tail.Options{})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
