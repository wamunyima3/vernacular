package sources

import (
	"slices"
	"strings"
	"testing"
)

func TestGatitosParse(t *testing.T) {
	tests := []struct {
		name      string
		parser    gatitos
		input     string
		want      []Record
		wantError bool
	}{
		{
			name:   "one record per target, text left as the source wrote it",
			parser: gatitos{file: "en_bem.jsonl"},
			input:  `{"src": "Hello", "trgs": ["Muli shani", " Halo "], "sl": "en", "tl": "bem"}`,
			want: []Record{
				{Text: "Muli shani", English: "Hello", Reference: "en_bem.jsonl:1"},
				{Text: " Halo ", English: "Hello", Reference: "en_bem.jsonl:1"},
			},
		},
		{
			name:   "skips blank targets",
			parser: gatitos{file: "en_bem.jsonl"},
			input:  "{\"src\": \"I\", \"trgs\": [\"ine\", \"  \"]}",
			want:   []Record{{Text: "ine", English: "I", Reference: "en_bem.jsonl:1"}},
		},
		{
			name:      "broken line reports file and line",
			parser:    gatitos{file: "en_bem.jsonl"},
			input:     "{\"src\": \"a\", \"trgs\": [\"b\"]}\nnot json",
			want:      []Record{{Text: "b", English: "a", Reference: "en_bem.jsonl:1"}},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []Record
			var gotErr error
			for rec, err := range tt.parser.Parse(strings.NewReader(tt.input)) {
				if err != nil {
					gotErr = err
					break
				}
				got = append(got, rec)
			}
			if (gotErr != nil) != tt.wantError {
				t.Fatalf("err = %v, wantError %v", gotErr, tt.wantError)
			}
			if gotErr != nil && !strings.Contains(gotErr.Error(), "en_bem.jsonl:2") {
				t.Errorf("error %q does not name the line", gotErr)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("got %+v\nwant %+v", got, tt.want)
			}
		})
	}
}
