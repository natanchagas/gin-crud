package ports

import (
	"context"

	"github.com/natanchagas/gin-crud/internal/core/domain"
)

//go:generate mockery --name RealStateRepository
type RealStateRepository interface {
	CreateRealState(ctx context.Context, realState domain.RealState) (int64, error)
}
