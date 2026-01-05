package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type EntityDecoder[T any] interface {
	ToEntity(*T) error
	Flush()
}

// EntityIterator is not safe for concurrent usage
type EntityIterator[T any, V EntityDecoder[T]] struct {
	cursor *mongo.Cursor
	v      V
}

func NewEntityStream[T any, V EntityDecoder[T]](
	cursor *mongo.Cursor,
	v V,
) *EntityIterator[T, V] {
	return &EntityIterator[T, V]{cursor: cursor, v: v}
}

func (s *EntityIterator[T, V]) Next(ctx context.Context) bool {
	return s.cursor.Next(ctx)
}

func (s *EntityIterator[T, V]) Decode(entity *T) error {
	if err := s.cursor.Decode(s.v); err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	if err := s.v.ToEntity(entity); err != nil {
		s.v.Flush()
		return fmt.Errorf("to entity: %w", err)
	}
	s.v.Flush()

	return nil
}

func (s *EntityIterator[T, V]) Close(ctx context.Context) error {
	if err := s.cursor.Close(ctx); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	return nil
}
