package phrases

import (
	"errors"
	"slices"
	"testing"
)

func TestPlaceholders(t *testing.T) {
	if got := Placeholders("Pay K{amount} to {name} by {date}, {name}"); !slices.Equal(got, []string{"amount", "name", "date"}) {
		t.Errorf("Placeholders = %v", got)
	}
	if got := Placeholders("{Name} {} {1x} plain"); got != nil {
		t.Errorf("Placeholders of malformed names = %v, want none", got)
	}
}

func TestCheckPlaceholders(t *testing.T) {
	tests := []struct {
		source, text string
		ok           bool
	}{
		{"Hello, {name}", "Mwapoleni, {name}", true},
		{"{count} of {total}", "{total}: {count}", true}, // order may change
		{"Hello, {name}", "Mwapoleni", false},
		{"Hello, {name}", "Mwapoleni, {mwine}", false},
		{"Hello", "Mwapoleni, {name}", false},
	}
	for _, tt := range tests {
		err := checkPlaceholders(tt.source, tt.text)
		if (err == nil) != tt.ok || (err != nil && !errors.Is(err, ErrPlaceholders)) {
			t.Errorf("checkPlaceholders(%q, %q) = %v, want ok %v", tt.source, tt.text, err, tt.ok)
		}
	}
}

func TestCheckSource(t *testing.T) {
	for source, ok := range map[string]bool{
		"Hello, {name}": true, "No placeholder": true, "Hello, {Name}": false, "Pay {first name}": false,
	} {
		if err := checkSource(source); (err == nil) != ok {
			t.Errorf("checkSource(%q) = %v, want ok %v", source, err, ok)
		}
	}
}
