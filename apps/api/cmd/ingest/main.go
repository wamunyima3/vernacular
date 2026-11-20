// Command ingest imports data into the lexicon.
//
//	ingest -feed gatitos-en-bem [-file local.jsonl]   a registered source, as `sourced` data
//	ingest -all                                        every feed in sources.Feeds
//	ingest -concepts data/concepts/x.csv               one of our own core concept lists
//	ingest -restem                                     recompute search stems after a stemmer change
package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/languages"
	"github.com/wamunyima3/vernacular/apps/api/internal/lexicon"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, "ingest:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	name := flag.String("feed", "", "feed to import: "+feedNames())
	file := flag.String("file", "", "read this local copy instead of downloading (checksum still verified)")
	concepts := flag.String("concepts", "", "CSV of core concepts (domain,gloss_en) to import instead of a feed")
	all := flag.Bool("all", false, "import every feed")
	restem := flag.Bool("restem", false, "recompute search stems for every active language")
	flag.Parse()

	pool, err := database.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer pool.Close()

	switch {
	case *restem:
		return restemActive(ctx, pool)
	case *concepts != "":
		return importConcepts(ctx, pool, *concepts)
	case *all:
		for _, feed := range sources.Feeds {
			if err := importFeed(ctx, pool, feed, ""); err != nil {
				return err
			}
		}
		return nil
	}
	feed, ok := sources.FeedByName(*name)
	if !ok {
		return fmt.Errorf("unknown feed %q; choose one of: %s", *name, feedNames())
	}
	return importFeed(ctx, pool, feed, *file)
}

func importFeed(ctx context.Context, db database.Beginner, feed sources.Feed, file string) error {
	data, err := load(ctx, feed.URL, file)
	if err != nil {
		return err
	}
	if sum := sha256.Sum256(data); hex.EncodeToString(sum[:]) != feed.SHA256 {
		return fmt.Errorf("checksum mismatch for %s: the file is not the pinned revision", feed.Name)
	}

	stats, err := sources.Ingest(ctx, db, feed, bytes.NewReader(data))
	if err != nil {
		return err
	}
	fmt.Printf("%s: %d records, %d new entries, %d new senses, %d new variants\n",
		feed.Name, stats.Records, stats.Entries, stats.Senses, stats.Variants)
	return nil
}

// restemActive recomputes search stems for every active language.
func restemActive(ctx context.Context, db database.Beginner) error {
	return database.InTx(ctx, db, "ingest:restem", func(tx pgx.Tx) error {
		langs, err := languages.NewStore(tx).List(ctx)
		if err != nil {
			return err
		}
		lex := lexicon.NewStore(tx)
		for _, l := range langs {
			if l.Status != string(gen.LanguageStatusActive) {
				continue
			}
			n, err := lex.Restem(ctx, l.Code)
			if err != nil {
				return err
			}
			fmt.Printf("%s: %d stems changed\n", l.Code, n)
		}
		return nil
	})
}

func importConcepts(ctx context.Context, db database.Beginner, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	added, err := sources.ImportConcepts(ctx, db, f)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	fmt.Printf("%s: %d new core concepts\n", path, added)
	return nil
}

func load(ctx context.Context, url, file string) ([]byte, error) {
	if file != "" {
		return os.ReadFile(file)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download %s: %s", url, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func feedNames() string {
	names := make([]string, len(sources.Feeds))
	for i, f := range sources.Feeds {
		names[i] = f.Name
	}
	return strings.Join(names, ", ")
}
