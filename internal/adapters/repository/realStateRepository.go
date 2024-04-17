package repository

import (
	"context"
	"database/sql"

	"github.com/natanchagas/gin-crud/internal/core/domain"
)

const (
	CreateRealState = `INSERT INTO real_states (real_state_registration, real_state_address, real_state_size, real_state_price, real_state_state) VALUES (?, ?, ?, ?, ?);`
	GetRealState    = `SELECT real_state_id, real_state_registration, real_state_address, real_state_size, real_state_price, real_state_state FROM real_states WHERE real_state_id = ?`
	UpdateRealState = `UPDATE real_states SET real_state_registration = ?, real_state_address = ?, real_state_size = ?, real_state_price = ?, real_state_state = ? WHERE real_state_id = ?`
	DeleteRealState = `DELETE FROM real_states WHERE real_state_id = ?`
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

func (r *realStateRepository) GetRealState(ctx context.Context, id uint64) (domain.RealState, error) {
	var realState domain.RealState

	row := r.db.QueryRowContext(ctx, GetRealState, id)
	if err := row.Scan(&realState.Id, &realState.Registration, &realState.Address, &realState.Size, &realState.Price, &realState.State); err != nil {
		return domain.RealState{}, err
	}

	return realState, nil
}

func (r *realStateRepository) UpdateRealState(ctx context.Context, realState domain.RealState, id uint64) (domain.RealState, error) {
	_, err := r.db.ExecContext(ctx, UpdateRealState, realState.Registration, realState.Address, realState.Size, realState.Price, realState.State, id)
	if err != nil {
		return domain.RealState{}, err
	}

	return realState, nil
}

func (r *realStateRepository) DeleteRealState(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx, DeleteRealState, id)
	if err != nil {
		return err
	}

	return nil
}
