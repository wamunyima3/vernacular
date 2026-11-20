package sources

// smolRevision pins the google/smol dataset (CC BY 4.0) to one commit.
const smolRevision = "fdaff3a1a019f89fa30a562a85d7c1d3e9150444"

func smolURL(path string) string {
	return "https://huggingface.co/datasets/google/smol/resolve/" + smolRevision + "/" + path
}

// Feeds lists every file the ingest command can import. Add a feed here only
// after its source has a register row with a licence that allows it.
//
// GATITOS bem_en.jsonl is not listed: it holds the same 4,450 pairs as en_bem.jsonl.
var Feeds = []Feed{
	{
		Name: "gatitos-en-bem", Source: "GATITOS", Language: "bem", Core: true,
		URL:    smolURL("gatitos/en_bem.jsonl"),
		SHA256: "a82cc6cfcd074c938756af242440b9f5364f4249dbdedd3709e87b332cab5310",
		Parser: gatitos{file: "gatitos/en_bem.jsonl"},
	},
}

// FeedByName finds a feed.
func FeedByName(name string) (Feed, bool) {
	for _, f := range Feeds {
		if f.Name == name {
			return f, true
		}
	}
	return Feed{}, false
}
