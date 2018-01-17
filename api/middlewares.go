package api

import (
	"github.com/baa-middleware/gzip"
	"github.com/baa-middleware/recovery"
	baa "gopkg.in/baa.v1"
)

func initMiddleware(b *baa.Baa) {
	// pannic recover
	b.Use(recovery.Recovery())

	// Gzip
	if baa.Env == baa.PROD {
		b.Use(gzip.Gzip(gzip.Options{CompressionLevel: 4}))
	}
}
