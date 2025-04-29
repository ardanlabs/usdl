package web

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey int

const (
	writerKey ctxKey = iota + 1
	traceIDKey
)

func SetTraceID(ctx context.Context, traceID uuid.UUID) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// GetTraceID returns the traceID for the request.
func GetTraceID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(traceIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}
	}

	return v
}

func setWriter(ctx context.Context, w http.ResponseWriter) context.Context {
	return context.WithValue(ctx, writerKey, w)
}

// GetWriter returns the underlying writer for the request.
func GetWriter(ctx context.Context) http.ResponseWriter {
	v, ok := ctx.Value(writerKey).(http.ResponseWriter)
	if !ok {
		return nil
	}

	return v
}
