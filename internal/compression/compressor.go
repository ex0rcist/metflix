package compression

import (
	"context"
	"net/http"

	"github.com/klauspost/compress/gzip"

	"github.com/ex0rcist/metflix/internal/logging"
)

// Struct to handle compression.
type Compressor struct {
	http.ResponseWriter

	context          context.Context
	encoder          *gzip.Writer
	supportedContent map[string]struct{}
}

// Constructor.
func NewCompressor(w http.ResponseWriter, ctx context.Context) *Compressor {
	supportedContent := map[string]struct{}{
		"application/json": {}, // {} uses no memory
		"text/html":        {},
	}

	return &Compressor{
		ResponseWriter:   w,
		context:          ctx,
		supportedContent: supportedContent,
	}
}

// Write body to response.
func (c *Compressor) Write(resp []byte) (int, error) {
	contentType := c.Header().Get("Content-Type")

	if _, ok := c.supportedContent[contentType]; !ok {
		logging.LogDebugCtx(c.context, "compression not supported for "+contentType)
		return c.ResponseWriter.Write(resp)
	}

	if c.encoder == nil {
		encoder, err := gzip.NewWriterLevel(c.ResponseWriter, gzip.BestSpeed)
		if err != nil {
			logging.LogErrorCtx(c.context, err)
			return c.ResponseWriter.Write(resp)
		}
		c.encoder = encoder
	}

	c.Header().Set("Content-Encoding", "gzip")

	return c.encoder.Write(resp)
}

// Close encoder.
func (c *Compressor) Close() {
	if c.encoder == nil {
		return
	}

	if err := c.encoder.Close(); err != nil {
		logging.LogErrorCtx(c.context, err, "error closing compressor encoder", err.Error())
	}

	c.encoder = nil
}
