package phrases

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
)

// Localisation files (PHR-3, US-12) key by the English source, so a pack
// drops into an i18n pipeline keyed by English. When a line has several
// translations, every format keeps the first.

// Map is a pack as English to translation, for the JSON format.
func Map(phrases []Phrase) map[string]string {
	out := make(map[string]string, len(phrases))
	for _, p := range firstPerSource(phrases) {
		out[p.SourceEn] = p.Text
	}
	return out
}

func firstPerSource(phrases []Phrase) []Phrase {
	seen := make(map[string]bool, len(phrases))
	out := phrases[:0:0]
	for _, p := range phrases {
		if !seen[p.SourceEn] {
			seen[p.SourceEn] = true
			out = append(out, p)
		}
	}
	return out
}

// PO is a gettext catalogue. Notes become translator comments, and
// the item's ID and trust status extracted comments.
func PO(pack, language string, phrases []Phrase) []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "# Vernacular phrase pack %q in %s.\n", pack, language)
	b.WriteString("msgid \"\"\nmsgstr \"\"\n")
	for _, h := range []string{
		"Project-Id-Version: vernacular " + pack, "Language-Team: Vernacular reviewers", "Language: " + language,
		"MIME-Version: 1.0", "Content-Type: text/plain; charset=UTF-8", "Content-Transfer-Encoding: 8bit",
		"X-Source-Language: en",
	} {
		b.WriteString(poQuote(h+"\n") + "\n")
	}
	b.WriteString("\n")
	for _, p := range firstPerSource(phrases) {
		for line := range strings.SplitSeq(p.Notes, "\n") {
			if line != "" {
				fmt.Fprintf(&b, "# %s\n", line)
			}
		}
		fmt.Fprintf(&b, "#. id: %s, status: %s\n", p.ID, p.Trust.Status)
		fmt.Fprintf(&b, "msgid %s\nmsgstr %s\n\n", poQuote(p.SourceEn), poQuote(p.Text))
	}
	return b.Bytes()
}

var poEscaper = strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`, "\t", `\t`, "\r", `\r`)

// poQuote is s as a PO string literal.
func poQuote(s string) string { return `"` + poEscaper.Replace(s) + `"` }

// XLIFF 2.0 (OASIS), one unit per phrase, keyed by the phrase's ID.
type xliffDoc struct {
	XMLName xml.Name  `xml:"urn:oasis:names:tc:xliff:document:2.0 xliff"`
	Version string    `xml:"version,attr"`
	SrcLang string    `xml:"srcLang,attr"`
	TrgLang string    `xml:"trgLang,attr"`
	File    xliffFile `xml:"file"`
}

type xliffFile struct {
	ID    string      `xml:"id,attr"`
	Units []xliffUnit `xml:"unit"`
}

type xliffUnit struct {
	ID      string       `xml:"id,attr"`
	Notes   *xliffNotes  `xml:"notes,omitempty"`
	Segment xliffSegment `xml:"segment"`
}

type xliffNotes struct {
	Note string `xml:"note"`
}

type xliffSegment struct {
	State  string `xml:"state,attr"`
	Source string `xml:"source"`
	Target string `xml:"target"`
}

// XLIFF is an XLIFF 2.0 document.
func XLIFF(pack, language string, phrases []Phrase) ([]byte, error) {
	doc := xliffDoc{Version: "2.0", SrcLang: "en", TrgLang: language, File: xliffFile{ID: pack}}
	for _, p := range firstPerSource(phrases) {
		u := xliffUnit{ID: p.ID.String(), Segment: xliffSegment{State: segmentState(p.Trust.Status), Source: p.SourceEn, Target: p.Text}}
		if p.Notes != "" {
			u.Notes = &xliffNotes{Note: p.Notes}
		}
		doc.File.Units = append(doc.File.Units, u)
	}
	var b bytes.Buffer
	b.WriteString(xml.Header)
	enc := xml.NewEncoder(&b)
	enc.Indent("", "  ")
	if err := enc.Encode(doc); err != nil {
		return nil, err
	}
	b.WriteString("\n")
	return b.Bytes(), nil
}

// segmentState maps trust to XLIFF's workflow: verified is final, reviewed
// is reviewed, anything less is translated but unchecked.
func segmentState(status string) string {
	switch status {
	case string(gen.TrustStatusVerified):
		return "final"
	case string(gen.TrustStatusReviewed):
		return "reviewed"
	}
	return "translated"
}
