package review_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestOnlyPromotionSetsStatus guards AGENTS.md rule 2: only the promotion
// function changes an item's status after it is created.
func TestOnlyPromotionSetsStatus(t *testing.T) {
	allowed := map[string]bool{
		filepath.Join("..", "review", "promotion.go"): true, // the caller
		filepath.Join("..", "items", "items.go"):      true, // the definition
	}
	err := filepath.WalkDir("..", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.Contains(string(src), ".SetStatus(") && !allowed[path] {
			t.Errorf("%s calls SetStatus; only review/promotion.go may", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	queries, err := filepath.Glob(filepath.Join("..", "..", "db", "queries", "*.sql"))
	if err != nil {
		t.Fatal(err)
	}
	setStatus := regexp.MustCompile(`(?i)update\s+items\s+set\s+status`)
	for _, q := range queries {
		src, err := os.ReadFile(q)
		if err != nil {
			t.Fatal(err)
		}
		if n := len(setStatus.FindAllIndex(src, -1)); n > 0 && (!strings.HasSuffix(q, "items.sql") || n != 1) {
			t.Errorf("%s updates items.status; only SetItemStatus in items.sql may", q)
		}
	}
}
