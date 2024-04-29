package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert/v2"

	"github.com/natanchagas/gin-crud/internal/adapters/repository"
	"github.com/natanchagas/gin-crud/internal/core/domain"
	"github.com/natanchagas/gin-crud/internal/pkg/customerrors"
)

func TestCreateRealState(t *testing.T) {
	type output struct {
		id  int64
		err error
	}

	testCases := []struct {
		name       string
		input      domain.RealState
		mocking    func(mock sqlmock.Sqlmock, realState domain.RealState) output
		assertions func(t *testing.T, actual, expected output)
	}{
		{
			name: "When real state is valid, should create a new real state",
			input: domain.RealState{
				Registration: 987654321,
				Address:      "456 Elm St",
				Size:         200,
				Price:        250000.50,
				State:        "CA",
			},
			mocking: func(mock sqlmock.Sqlmock, realState domain.RealState) output {
				mock.
					ExpectExec("INSERT INTO real_states").
					WithArgs(realState.Registration, realState.Address, realState.Size, realState.Price, realState.State).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return output{
					id:  1,
					err: nil,
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, actual, expected)
			},
		},
		{
			name: "When real state is valid, but some error happens on last id request, should return error",
			input: domain.RealState{
				Registration: 987654321,
				Address:      "456 Elm St",
				Size:         200,
				Price:        250000.50,
				State:        "CA",
			},
			mocking: func(mock sqlmock.Sqlmock, realState domain.RealState) output {
				mock.
					ExpectExec("INSERT INTO real_states").
					WithArgs(realState.Registration, realState.Address, realState.Size, realState.Price, realState.State).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("unexpected error")))

				return output{
					id:  -1,
					err: customerrors.Internal,
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, actual, expected)
			},
		},
		{
			name: "When real state is valid, but insert fails, should return error",
			input: domain.RealState{
				Registration: 987654321,
				Address:      "456 Elm St",
				Size:         200,
				Price:        250000.50,
				State:        "CA",
			},
			mocking: func(mock sqlmock.Sqlmock, realState domain.RealState) output {
				mock.
					ExpectExec("INSERT INTO real_states").
					WithArgs(realState.Registration, realState.Address, realState.Size, realState.Price, realState.State).
					WillReturnError(errors.New("insert error"))

				return output{
					id:  -1,
					err: customerrors.Internal,
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, actual, expected)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			expected := tc.mocking(mock, tc.input)

			r := repository.NewRealStateRepository(db)
			var actual output

			actual.id, actual.err = r.CreateRealState(ctx, tc.input)

			tc.assertions(t, expected, actual)

		})
	}
}

func TestGetRealState(t *testing.T) {
	type output struct {
		realState domain.RealState
		err       error
	}

	testCases := []struct {
		name       string
		input      uint64
		mocking    func(mock sqlmock.Sqlmock, id uint64) output
		assertions func(t *testing.T, actual, expected output)
	}{
		{
			name:  "When real state exists, should return real state",
			input: 1,
			mocking: func(mock sqlmock.Sqlmock, id uint64) output {
				mock.
					ExpectQuery(`SELECT real_state_id, real_state_registration, real_state_address, real_state_size, real_state_price, real_state_state FROM real_states WHERE real_state_id = ?`).
					WithArgs(id).
					WillReturnRows(sqlmock.NewRows([]string{"real_state_id", "real_state_registration", "real_state_address", "real_state_size", "real_state_price", "real_state_state"}).
						AddRow(1, 987654321, "456 Elm St", 200, 250000.50, "CA"))

				return output{
					realState: domain.RealState{
						Id:           1,
						Registration: 987654321,
						Address:      "456 Elm St",
						Size:         200,
						Price:        250000.50,
						State:        "CA",
					},
					err: nil,
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, actual, expected)
			},
		},
		{
			name:  "When real state does not exists, should return error",
			input: 1,
			mocking: func(mock sqlmock.Sqlmock, id uint64) output {
				mock.
					ExpectQuery(`SELECT real_state_id, real_state_registration, real_state_address, real_state_size, real_state_price, real_state_state FROM real_states WHERE real_state_id = ?`).
					WithArgs(id).
					WillReturnError(sql.ErrNoRows)

				return output{
					realState: domain.RealState{},
					err:       customerrors.NotFound,
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, actual, expected)
			},
		},
		{
			name:  "When real state exists, but something fails should return error",
			input: 1,
			mocking: func(mock sqlmock.Sqlmock, id uint64) output {
				mock.
					ExpectQuery(`SELECT real_state_id, real_state_registration, real_state_address, real_state_size, real_state_price, real_state_state FROM real_states WHERE real_state_id = ?`).
					WithArgs(id).
					WillReturnError(sql.ErrConnDone)

				return output{
					realState: domain.RealState{},
					err:       customerrors.Internal,
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, actual, expected)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			expected := tc.mocking(mock, tc.input)

			r := repository.NewRealStateRepository(db)
			var actual output
			actual.realState, actual.err = r.GetRealState(ctx, tc.input)

			tc.assertions(t, expected, actual)
		})
	}
}

func TestUpdateRealState(t *testing.T) {
	type output struct {
		realState domain.RealState
		err       error
	}

	type input struct {
		realState domain.RealState
		id        uint64
	}

	testCases := []struct {
		name       string
		input      input
		mocking    func(mock sqlmock.Sqlmock, in input) output
		assertions func(t *testing.T, actual, expected output)
	}{
		{
			name: "When real state exists, should update real state",
			input: input{
				realState: domain.RealState{
					Id:           1,
					Registration: 987654321,
					Address:      "456 Elm St",
					Size:         200,
					Price:        275000.00,
					State:        "CA",
				},
				id: 1,
			},
			mocking: func(mock sqlmock.Sqlmock, in input) output {
				mock.
					ExpectExec(`UPDATE real_state`).
					WithArgs(in.realState.Registration, in.realState.Address, in.realState.Size, in.realState.Price, in.realState.State, in.id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return output{
					realState: in.realState,
					err:       nil,
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, actual, expected)
			},
		},
		{
			name: "When real state exists, but update fails, should return error",
			input: input{
				realState: domain.RealState{
					Id:           1,
					Registration: 987654321,
					Address:      "456 Elm St",
					Size:         200,
					Price:        275000.00,
					State:        "CA",
				},
				id: 1,
			},
			mocking: func(mock sqlmock.Sqlmock, in input) output {
				mock.
					ExpectExec(`UPDATE real_state`).
					WithArgs(in.realState.Registration, in.realState.Address, in.realState.Size, in.realState.Price, in.realState.State, in.id).
					WillReturnError(errors.New("update failed"))

				return output{
					realState: domain.RealState{},
					err:       customerrors.Internal,
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, actual, expected)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			expected := tc.mocking(mock, tc.input)

			r := repository.NewRealStateRepository(db)
			var actual output
			actual.realState, actual.err = r.UpdateRealState(ctx, tc.input.realState, tc.input.id)

			tc.assertions(t, expected, actual)
		})
	}
}

func TestDeleteRealState(t *testing.T) {
	testCases := []struct {
		name       string
		input      uint64
		mocking    func(mock sqlmock.Sqlmock, id uint64) error
		assertions func(t *testing.T, actual, expected error)
	}{
		{
			name:  "When real state exists, should delete real state",
			input: 1,
			mocking: func(mock sqlmock.Sqlmock, id uint64) error {
				mock.
					ExpectExec(`DELETE FROM real_states`).
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return nil
			},
			assertions: func(t *testing.T, actual, expected error) {
				assert.Equal(t, actual, expected)
			},
		},
		{
			name:  "When real state exists, but delete fails, should return error",
			input: 1,
			mocking: func(mock sqlmock.Sqlmock, id uint64) error {
				mock.
					ExpectExec(`DELETE FROM real_states`).
					WithArgs(id).
					WillReturnError(errors.New("delete failed"))

				return customerrors.Internal
			},
			assertions: func(t *testing.T, actual, expected error) {
				assert.Equal(t, actual, expected)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			expected := tc.mocking(mock, tc.input)

			r := repository.NewRealStateRepository(db)
			actual := r.DeleteRealState(ctx, tc.input)

			tc.assertions(t, expected, actual)
		})
	}
}
