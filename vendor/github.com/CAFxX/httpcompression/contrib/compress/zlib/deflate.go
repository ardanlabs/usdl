package zlib

import (
	"compress/zlib"
	"fmt"
	"io"
	"sync"

	"github.com/CAFxX/httpcompression/contrib/internal/utils"
)

const (
	Encoding           = "deflate"
	DefaultCompression = zlib.DefaultCompression
)

type Options struct {
	Level      int
	Dictionary []byte
}

type compressor struct {
	pool sync.Pool
	opt  Options
}

func New(opt Options) (*compressor, error) {
	tw, err := zlib.NewWriterLevelDict(io.Discard, opt.Level, opt.Dictionary)
	if err != nil {
		return nil, err
	}
	err = utils.CheckWriter(tw)
	if err != nil {
		return nil, fmt.Errorf("deflate: writer initialization: %w", err)
	}

	c := &compressor{opt: opt}
	return c, nil
}

func (c *compressor) Get(w io.Writer) io.WriteCloser {
	if gw, ok := c.pool.Get().(*deflateWriter); ok {
		gw.Reset(w)
		return gw
	}
	gw, err := zlib.NewWriterLevelDict(w, c.opt.Level, c.opt.Dictionary)
	if err != nil {
		return utils.ErrorWriteCloser{Err: err}
	}
	return &deflateWriter{
		Writer: gw,
		c:      c,
	}
}

type deflateWriter struct {
	*zlib.Writer
	c *compressor
}

func (w *deflateWriter) Close() error {
	err := w.Writer.Close()
	w.Reset(nil)
	w.c.pool.Put(w)
	return err
}
