package compression

import (
	"context"
	"net/http"

	"github.com/klauspost/compress/gzip"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
)

type Decompressor struct {
	request            *http.Request
	reader             *gzip.Reader
	context            context.Context
	supportedEncodings map[string]struct{}
}

func NewDecompressor(req *http.Request, ctx context.Context) *Decompressor {
	supportedEncodings := map[string]struct{}{
		"gzip": {}, // {} uses no memory
	}

	return &Decompressor{
		request:            req,
		context:            ctx,
		supportedEncodings: supportedEncodings,
	}
}

func (d *Decompressor) Decompress() error {
	encoding := d.request.Header.Get("Content-Encoding")

	if len(encoding) == 0 {
		logging.LogDebugCtx(d.context, "no encoding provided")
		return nil
	}

	logging.LogDebugCtx(d.context, "got request compressed with "+encoding)

	if _, ok := d.supportedEncodings[encoding]; !ok {
		err := entities.ErrEncodingUnsupported
		logging.LogErrorCtx(d.context, err, "decoding not supported for "+encoding)

		return err
	}

	if d.reader == nil {
		reader, err := gzip.NewReader(d.request.Body)
		if err != nil {
			err := entities.ErrEncodingInternal
			logging.LogErrorCtx(d.context, err, "failed to create gzip reader: "+err.Error())

			return err
		}

		d.reader = reader
	}

	d.request.Body = d.reader

	return nil
}

func (d *Decompressor) Close() {
	if d.reader == nil {
		return
	}

	if err := d.reader.Close(); err != nil {
		logging.LogErrorCtx(d.context, err, "error closing decompressor reader", err.Error())
	}

	d.reader = nil
}
