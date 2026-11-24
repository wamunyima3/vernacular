package server_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
)

// TestContractHasNoEmptyKeys: a comma inside a YAML flow mapping, such as
// `{ type: string, description: Calls today, UTC }`, silently turns the text
// after it into a key with no value. The contract never needs a null, so any
// null is such a slip.
func TestContractHasNoEmptyKeys(t *testing.T) {
	raw, err := gen.PathToRawSpec("spec")["spec"]()
	if err != nil {
		t.Fatal(err)
	}
	var doc any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	var walk func(v any, path string)
	walk = func(v any, path string) {
		switch v := v.(type) {
		case nil:
			t.Errorf("%s has no value: quote the text before it, or write the mapping as a block", path)
		case map[string]any:
			for k, x := range v {
				walk(x, path+"/"+k)
			}
		case []any:
			for i, x := range v {
				walk(x, fmt.Sprintf("%s/%d", path, i))
			}
		}
	}
	walk(doc, "")
}
