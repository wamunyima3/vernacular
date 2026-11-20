// Package blob opens the object store for permission documents, audio and
// release exports (architecture §9: S3-compatible only). It uses gocloud.dev,
// so the same code runs against S3, any S3-compatible server, a local folder
// or memory, chosen by URL:
//
//	s3://bucket?endpoint=http://host:8333&use_path_style=true&region=us-east-1
//	file:///abs/path?create_dir=true
//	mem://
package blob

import (
	"context"
	"fmt"

	"gocloud.dev/blob"

	_ "gocloud.dev/blob/fileblob" // file:// URLs, for development
	_ "gocloud.dev/blob/memblob"  // mem:// URLs, for tests
	_ "gocloud.dev/blob/s3blob"   // s3:// URLs, for production
)

// Bucket is a gocloud bucket. Callers use its NewReader, WriteAll, Exists and
// Delete methods directly.
type Bucket = blob.Bucket

// Open opens the bucket at url.
func Open(ctx context.Context, url string) (*Bucket, error) {
	b, err := blob.OpenBucket(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("open blob bucket: %w", err)
	}
	return b, nil
}
