package sources

// smolRevision pins the google/smol dataset (CC BY 4.0) to one commit.
const smolRevision = "fdaff3a1a019f89fa30a562a85d7c1d3e9150444"

// gatitosFeed builds a feed for one GATITOS en_xx file of the pinned revision.
func gatitosFeed(name, language, path, sha256 string) Feed {
	return Feed{
		Name: name, Source: "GATITOS", Language: language, Core: true,
		URL:    "https://huggingface.co/datasets/google/smol/resolve/" + smolRevision + "/" + path,
		SHA256: sha256,
		Parser: gatitos{file: path},
	}
}

// Feeds lists every file the ingest command can import. Add a feed here only
// after its source has a register row with a licence that allows it.
//
// GATITOS bem_en.jsonl is not listed: it holds the same 4,450 pairs as en_bem.jsonl.
var Feeds = []Feed{
	gatitosFeed("gatitos-en-bem", "bem", "gatitos/en_bem.jsonl",
		"a82cc6cfcd074c938756af242440b9f5364f4249dbdedd3709e87b332cab5310"),
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
