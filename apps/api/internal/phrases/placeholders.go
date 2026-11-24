package phrases

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// Placeholders (PHR-2, US-32) are names in braces, such as {name} or
// {amount}, that an app fills in. A translation keeps every placeholder of
// its English line, and its notes say how the words around one change with
// what fills it, such as noun class agreement.
var placeholderPattern = regexp.MustCompile(`\{([a-z][a-z0-9_]*)\}`)

// ErrPlaceholders means a translation does not keep its line's placeholders.
var ErrPlaceholders = errors.New("a translation must keep the placeholders of its English line")

// PlaceholderMessage is err's reason, written for a person.
func PlaceholderMessage(err error) string {
	msg := err.Error()
	return strings.ToUpper(msg[:1]) + msg[1:]
}

// malformed finds braces that are not a placeholder, such as {Name} or
// {first name}, which apps would print as written.
var malformed = regexp.MustCompile(`\{[^}]*\}`)

// checkSource refuses an English line with a malformed placeholder.
func checkSource(source string) error {
	for _, m := range malformed.FindAllString(source, -1) {
		if !placeholderPattern.MatchString(m) {
			return fmt.Errorf("%q: placeholder %s must be lower-case letters, digits and _", source, m)
		}
	}
	return nil
}

// Placeholders lists the placeholder names in s, once each, in order.
func Placeholders(s string) []string {
	var out []string
	for _, m := range placeholderPattern.FindAllStringSubmatch(s, -1) {
		if !slices.Contains(out, m[1]) {
			out = append(out, m[1])
		}
	}
	return out
}

// checkPlaceholders refuses text whose placeholders are not exactly those
// of the English source, in any order.
func checkPlaceholders(source, text string) error {
	want, got := Placeholders(source), Placeholders(text)
	slices.Sort(want)
	slices.Sort(got)
	if slices.Equal(want, got) {
		return nil
	}
	if len(want) == 0 {
		return fmt.Errorf("%w: this line has none", ErrPlaceholders)
	}
	return fmt.Errorf("%w: keep {%s}", ErrPlaceholders, strings.Join(want, "}, {"))
}
