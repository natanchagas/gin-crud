package repository

import (
	"context"
	"database/sql"

	"github.com/natanchagas/gin-crud/internal/core/domain"
)

const (
	CreateRealState = `INSERT INTO real_states (real_state_registration, real_state_address, real_state_size, real_state_price, real_state_state) VALUES (?, ?, ?, ?, ?);`
)

type realStateRepository struct {
	db *sql.DB
}

func NewRealStateRepository(db *sql.DB) *realStateRepository {
	return &realStateRepository{
		db: db,
	}
}

func (r *realStateRepository) CreateRealState(ctx context.Context, realState domain.RealState) (int64, error) {
	res, err := r.db.ExecContext(ctx, CreateRealState, realState.Registration, realState.Address, realState.Size, realState.Price, realState.State)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}
