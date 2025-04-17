package zstd

import (
	"fmt"
	"io"
	"sync"

	"github.com/CAFxX/httpcompression/contrib/internal/utils"
	"github.com/klauspost/compress/zstd"
)

const (
	Encoding           = "zstd"
	DefaultCompression = zstd.SpeedDefault
)

type compressor struct {
	pool sync.Pool
	opts []zstd.EOption
}

func New(opts ...zstd.EOption) (c *compressor, err error) {
	defer func() {
		if r := recover(); r != nil {
			c, err = nil, fmt.Errorf("panic: %v", r)
		}
	}()

	opts = append([]zstd.EOption(nil), opts...)

	tw, err := zstd.NewWriter(io.Discard, opts...)
	if err != nil {
		return nil, err
	}
	if err := utils.CheckWriter(tw); err != nil {
		return nil, fmt.Errorf("zstd: writer initialization: %w", err)
	}

	c = &compressor{opts: opts}
	return c, nil
}

func (c *compressor) Get(w io.Writer) io.WriteCloser {
	if gw, ok := c.pool.Get().(*zstdWriter); ok {
		gw.Reset(w)
		return gw
	}
	gw, err := zstd.NewWriter(w, c.opts...)
	if err != nil {
		return utils.ErrorWriteCloser{Err: err}
	}
	return &zstdWriter{
		Encoder: gw,
		c:       c,
	}
}

type zstdWriter struct {
	*zstd.Encoder
	c *compressor
}

func (w *zstdWriter) Close() error {
	err := w.Encoder.Close()
	w.Reset(nil)
	w.c.pool.Put(w)
	return err
}
