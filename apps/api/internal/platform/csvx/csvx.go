// Package csvx reads our own CSV lists (concepts, phrases) row by row.
package csvx

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"slices"
)

// Each checks the header, then calls fn for every row. An error from fn is
// returned with the row's line number.
func Each(r io.Reader, header []string, fn func(rec []string) error) error {
	rows := csv.NewReader(r)
	got, err := rows.Read()
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}
	if !slices.Equal(got, header) {
		return fmt.Errorf("header = %v, want %v", got, header)
	}
	for {
		rec, err := rows.Read()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if err := fn(rec); err != nil {
			line, _ := rows.FieldPos(0)
			return fmt.Errorf("line %d: %w", line, err)
		}
	}
}
