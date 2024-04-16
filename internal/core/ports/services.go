package ports

import (
	"context"

	"github.com/natanchagas/gin-crud/internal/core/domain"
)

//go:generate mockery --name RealStateService
type RealStateService interface {
	Create(ctx context.Context, realState domain.RealState) (domain.RealState, error)
}
