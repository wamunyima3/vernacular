// Command ingest imports a registered source file into the lexicon as
// `sourced` data. Usage: ingest -feed gatitos-en-bem [-file local.jsonl]
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
	flag.Parse()

	feed, ok := sources.FeedByName(*name)
	if !ok {
		return fmt.Errorf("unknown feed %q; choose one of: %s", *name, feedNames())
	}
	data, err := load(ctx, feed.URL, *file)
	if err != nil {
		return err
	}
	if sum := sha256.Sum256(data); hex.EncodeToString(sum[:]) != feed.SHA256 {
		return fmt.Errorf("checksum mismatch for %s: the file is not the pinned revision", feed.Name)
	}

	pool, err := database.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer pool.Close()

	stats, err := sources.Ingest(ctx, pool, feed, bytes.NewReader(data))
	if err != nil {
		return err
	}
	fmt.Printf("%s: %d records, %d new entries, %d new senses\n", feed.Name, stats.Records, stats.Entries, stats.Senses)
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
