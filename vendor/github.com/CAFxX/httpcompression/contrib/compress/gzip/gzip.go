package gzip

import (
	"compress/gzip"
	"fmt"
	"io"
	"sync"

	"github.com/CAFxX/httpcompression/contrib/internal/utils"
)

const (
	Encoding           = "gzip"
	DefaultCompression = gzip.DefaultCompression
)

type Options struct {
	Level int
}

type compressor struct {
	pool sync.Pool
	opt  Options
}

func New(opt Options) (*compressor, error) {
	tw, err := gzip.NewWriterLevel(io.Discard, opt.Level)
	if err != nil {
		return nil, err
	}
	err = utils.CheckWriter(tw)
	if err != nil {
		return nil, fmt.Errorf("gzip: writer initialization: %w", err)
	}

	c := &compressor{opt: opt}
	return c, nil
}

func (c *compressor) Get(w io.Writer) io.WriteCloser {
	if gw, ok := c.pool.Get().(*gzipWriter); ok {
		gw.Reset(w)
		return gw
	}
	gw, err := gzip.NewWriterLevel(w, c.opt.Level)
	if err != nil {
		return utils.ErrorWriteCloser{Err: err}
	}
	return &gzipWriter{
		Writer: gw,
		c:      c,
	}
}

type gzipWriter struct {
	*gzip.Writer
	c *compressor
}

func (w *gzipWriter) Close() error {
	err := w.Writer.Close()
	w.Reset(nil)
	w.c.pool.Put(w)
	return err
}
