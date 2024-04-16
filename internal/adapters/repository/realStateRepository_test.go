package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert/v2"

	"github.com/natanchagas/gin-crud/internal/adapters/repository"
	"github.com/natanchagas/gin-crud/internal/core/domain"
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
					err: errors.New("unexpected error"),
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
					err: errors.New("insert error"),
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
