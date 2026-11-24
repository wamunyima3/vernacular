// Search the Bemba dictionary for "amenshi", unchecked words included.
// Run with: go run search.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type results struct {
	Data []struct {
		Entry struct {
			Headword string `json:"headword"`
			Senses   []struct {
				Gloss string `json:"gloss"`
				Trust struct {
					Status string `json:"status"`
				} `json:"trust"`
			} `json:"senses"`
		} `json:"entry"`
	} `json:"data"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func main() {
	query := url.Values{"language": {"bem"}, "q": {"amenshi"}, "min_status": {"machine"}}
	req, err := http.NewRequest(http.MethodGet, os.Getenv("VERNACULAR_URL")+"/v1/entries/search?"+query.Encode(), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-API-Key", os.Getenv("VERNACULAR_KEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	var body results
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("%d %s: %s", res.StatusCode, body.Code, body.Message)
	}

	for _, hit := range body.Data {
		meanings := make([]string, len(hit.Entry.Senses))
		for i, s := range hit.Entry.Senses {
			meanings[i] = fmt.Sprintf("%s (%s)", s.Gloss, s.Trust.Status)
		}
		fmt.Printf("%s: %s\n", hit.Entry.Headword, strings.Join(meanings, ", "))
	}
}
