package ports

import (
	"context"

	"github.com/natanchagas/gin-crud/internal/core/domain"
)

//go:generate mockery --name RealStateService
type RealStateService interface {
	Create(ctx context.Context, realState domain.RealState) (domain.RealState, error)
	Get(ctx context.Context, id uint64) (domain.RealState, error)
	Update(ctx context.Context, realState domain.RealState, id uint64) (domain.RealState, error)
	Delete(ctx context.Context, id uint64) error
}
