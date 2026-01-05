package stream

import (
	"context"
)

type Stream[T any] interface {
	Next(ctx context.Context) bool
	Decode(e *T) error
	Close(ctx context.Context) error
}
