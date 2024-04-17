package ports

import (
	"context"

	"github.com/natanchagas/gin-crud/internal/core/domain"
)

//go:generate mockery --name RealStateRepository
type RealStateRepository interface {
	CreateRealState(ctx context.Context, realState domain.RealState) (int64, error)
	GetRealState(ctx context.Context, id uint64) (domain.RealState, error)
	UpdateRealState(ctx context.Context, realState domain.RealState, id uint64) (domain.RealState, error)
	DeleteRealState(ctx context.Context, id uint64) error
}
