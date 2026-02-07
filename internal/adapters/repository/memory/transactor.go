package memory

import (
	"context"

	"github.com/shanth1/gotrace/internal/core/ports"
)

type MemoryTransactor struct{}

func NewMemoryTransactor() ports.Transactor {
	return &MemoryTransactor{}
}

func (t *MemoryTransactor) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	return tFunc(ctx)
}
