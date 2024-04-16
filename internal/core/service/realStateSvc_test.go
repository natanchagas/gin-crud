package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/natanchagas/gin-crud/internal/core/domain"
	"github.com/natanchagas/gin-crud/internal/core/service"
	"github.com/natanchagas/gin-crud/internal/mocks"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	type output struct {
		realState domain.RealState
		err       error
	}

	testCases := []struct {
		name      string
		input     domain.RealState
		mocking   func(mock *mocks.RealStateRepository, realState domain.RealState) output
		assertion func(t *testing.T, actual, expected output)
	}{
		{
			name: "When real state is valid and repository returns success, should return real state with id",
			input: domain.RealState{
				Registration: 987654321,
				Address:      "456 Elm St",
				Size:         200,
				Price:        250000.50,
				State:        "CA",
			},
			mocking: func(m *mocks.RealStateRepository, realState domain.RealState) output {

				m.
					On("CreateRealState", mock.AnythingOfType("context.backgroundCtx"), realState).
					Return(int64(1), nil)

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
			assertion: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "When real state is valid, but repository fails, should return real state with id",
			input: domain.RealState{
				Registration: 987654321,
				Address:      "456 Elm St",
				Size:         200,
				Price:        250000.50,
				State:        "CA",
			},
			mocking: func(m *mocks.RealStateRepository, realState domain.RealState) output {

				m.
					On("CreateRealState", mock.AnythingOfType("context.backgroundCtx"), realState).
					Return(int64(-1), errors.New("failed to create"))

				return output{
					realState: domain.RealState{},
					err:       errors.New("failed to create"),
				}
			},
			assertion: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected, actual)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			r := mocks.NewRealStateRepository(t)
			s := service.NewRealStateService(r)

			expected := tc.mocking(r, tc.input)

			var actual output
			actual.realState, actual.err = s.Create(ctx, tc.input)

			tc.assertion(t, actual, expected)
		})
	}
}
