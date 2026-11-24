package releases

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/releases/dbgen"
)

// table is one CSV file of the open dataset. Queries run pinned, so
// unqualified item tables are the release's frozen copies (ADR-0012).
type table struct {
	file, about, query string
}

var tables = []table{
	{
		"entries.csv", "Dictionary entries: headword, part of speech, noun class, plural, root, stem, variety and register",
		`SELECT id, language_code AS language, headword, part_of_speech, noun_class, plural, root, stem, variety, register
		 FROM entries ORDER BY language_code, headword, id`,
	},
	{
		"variants.csv", "Other spellings of an entry",
		`SELECT entry_id, spelling, note FROM variants ORDER BY entry_id, spelling`,
	},
	{
		"senses.csv", "Meanings of an entry, in order, each linked to a concept",
		`SELECT id, entry_id, position, gloss, definition, notes, concept_id FROM senses ORDER BY entry_id, position, id`,
	},
	{
		"examples.csv", "Example sentences of a sense, with an English translation",
		`SELECT id, sense_id, text, translation_en FROM examples ORDER BY sense_id, id`,
	},
	{
		"phrases.csv", "Phrase pack lines: the English source and its translation",
		`SELECT id, language_code AS language, pack, source_en, text, notes, concept_id
		 FROM phrases ORDER BY language_code, pack, source_en, id`,
	},
	{
		"concepts.csv", "Language-independent meanings that senses and phrases express",
		`SELECT id, gloss_en, definition, domain FROM concepts
		 WHERE id IN (SELECT concept_id FROM senses UNION SELECT concept_id FROM phrases)
		 ORDER BY gloss_en, id`,
	},
	{
		"item_sources.csv", "Where each entry, sense, example and phrase comes from",
		`SELECT x.item_id, s.name AS source, x.reference, x.model
		 FROM item_sources x JOIN sources s ON s.id = x.source_id ORDER BY x.item_id, s.name, x.reference`,
	},
	{
		"sources.csv", "Each cited source, its rights holder, licence and allowed uses",
		`SELECT name, kind, rights_holder, licence, url, array_to_string(allowed_uses, ' ') AS allowed_uses
		 FROM sources WHERE id IN (SELECT source_id FROM item_sources) ORDER BY name`,
	},
}

// Export returns the named release, or the newest for `latest`, as a zip of
// CSV files with release.json and a README (NFR-9, US-17), and the release's
// name. It is built in memory: a release is small, and the HTTP layer needs
// the whole body for its ETag anyway.
func (s *Service) Export(ctx context.Context, name string) (string, []byte, error) {
	var w bytes.Buffer
	var release string
	err := s.pinned(ctx, name, func(tx pgx.Tx, row dbgen.Release) error {
		release = row.Name
		r, err := fromRow(row)
		if err != nil {
			return err
		}
		z := zip.NewWriter(&w)
		if err := writeFile(z, "README.md", []byte(readme(r))); err != nil {
			return err
		}
		meta, err := json.MarshalIndent(toMeta(r), "", "  ")
		if err != nil {
			return err
		}
		if err := writeFile(z, "release.json", meta); err != nil {
			return err
		}
		for _, t := range tables {
			f, err := z.Create(t.file)
			if err != nil {
				return err
			}
			// Postgres writes the CSV, quoting and all.
			if _, err := tx.Conn().PgConn().CopyTo(ctx, f, "COPY ("+t.query+") TO STDOUT WITH (FORMAT csv, HEADER)"); err != nil {
				return fmt.Errorf("export %s: %w", t.file, err)
			}
		}
		return z.Close()
	})
	return release, w.Bytes(), err
}

func writeFile(z *zip.Writer, name string, body []byte) error {
	f, err := z.Create(name)
	if err != nil {
		return err
	}
	_, err = f.Write(body)
	return err
}

// meta is release.json: the release as GET /v1/releases gives it, and
// what each file holds, for programs.
type meta struct {
	gen.Release
	Files map[string]string `json:"files"`
}

func toMeta(r Release) meta {
	files := make(map[string]string, len(tables))
	for _, t := range tables {
		files[t.file] = t.about
	}
	return meta{Release: toAPI(r), Files: files}
}

// readme describes the release for people.
func readme(r Release) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Vernacular dataset release %s\n\n", r.Name)
	fmt.Fprintf(&b, "Published %s. Every item here was verified by two independent reviewers.\n\n", r.PublishedAt.UTC().Format("2 January 2006"))
	if r.Notes != "" {
		fmt.Fprintf(&b, "%s\n\n", r.Notes)
	}
	fmt.Fprintf(&b, "Since the release before: %d added, %d changed, %d removed.\n\n", r.Changes.Added, r.Changes.Changed, r.Changes.Removed)
	b.WriteString("## Items per language\n\n| Language | Entries | Senses | Examples | Phrases |\n|---|---:|---:|---:|---:|\n")
	codes := make([]string, 0, len(r.Languages))
	for code := range r.Languages {
		codes = append(codes, code)
	}
	slices.Sort(codes)
	for _, code := range codes {
		k := r.Languages[code]
		fmt.Fprintf(&b, "| %s | %d | %d | %d | %d |\n", code, k["entry"], k["sense"], k["example"], k["phrase"])
	}
	b.WriteString("\n## Files\n\nUTF-8 CSV with a header row. IDs are UUIDs and join the files together.\n\n")
	for _, t := range tables {
		fmt.Fprintf(&b, "- `%s`: %s\n", t.file, t.about)
	}
	b.WriteString("- `release.json`: this summary, for programs\n\n")
	b.WriteString("## Licence and credit\n\n" +
		"Each item's sources are in `item_sources.csv`, and each source's licence in `sources.csv`. " +
		"Credit sources the way their licences ask. See the licences page of the documentation for the terms of the dataset as a whole.\n")
	return b.String()
}
