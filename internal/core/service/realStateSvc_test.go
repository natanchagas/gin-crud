package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/natanchagas/gin-crud/internal/core/domain"
	"github.com/natanchagas/gin-crud/internal/core/service"
	"github.com/natanchagas/gin-crud/internal/mocks"
	"github.com/stretchr/testify/assert"
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

func TestGet(t *testing.T) {
	type output struct {
		realState domain.RealState
		err       error
	}

	testCases := []struct {
		name      string
		input     uint64
		mocking   func(mock *mocks.RealStateRepository, id uint64) output
		assertion func(t *testing.T, actual, expected output)
	}{
		{
			name:  "When real state exists, should return it",
			input: 1,
			mocking: func(m *mocks.RealStateRepository, id uint64) output {
				m.
					On("GetRealState", mock.AnythingOfType("context.backgroundCtx"), id).
					Return(
						domain.RealState{
							Id:           1,
							Registration: 987654321,
							Address:      "456 Elm St",
							Size:         200,
							Price:        250000.50,
							State:        "CA",
						},
						nil,
					)

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
			name:  "When real state does not exists, should return error",
			input: 1,
			mocking: func(m *mocks.RealStateRepository, id uint64) output {
				m.
					On("GetRealState", mock.AnythingOfType("context.backgroundCtx"), id).
					Return(
						domain.RealState{},
						errors.New("real state not found"),
					)

				return output{
					realState: domain.RealState{},
					err:       errors.New("real state not found"),
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
			actual.realState, actual.err = s.Get(ctx, tc.input)

			tc.assertion(t, actual, expected)
		})
	}
}

func TestUpdate(t *testing.T) {
	type input struct {
		realState domain.RealState
		id        uint64
	}

	type output struct {
		realState domain.RealState
		err       error
	}

	testCases := []struct {
		name      string
		input     input
		mocking   func(m *mocks.RealStateRepository, in input) output
		assertion func(t *testing.T, actual, expected output)
	}{
		{
			name: "When real state is updated, should return real state",
			input: input{
				realState: domain.RealState{
					Registration: 987654321,
					Address:      "456 Elm St",
					Size:         200,
					Price:        275000.00,
					State:        "CA",
				},
				id: 1,
			},
			mocking: func(m *mocks.RealStateRepository, in input) output {
				m.
					On("GetRealState", mock.AnythingOfType("context.backgroundCtx"), in.id).
					Return(
						in.realState,
						nil,
					)

				m.
					On("UpdateRealState", mock.AnythingOfType("context.backgroundCtx"), in.realState, in.id).
					Return(
						in.realState,
						nil,
					)

				return output{
					realState: in.realState,
					err:       nil,
				}
			},
			assertion: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "When real state update fails, should return error",
			input: input{
				realState: domain.RealState{
					Registration: 987654321,
					Address:      "456 Elm St",
					Size:         200,
					Price:        275000.00,
					State:        "CA",
				},
				id: 1,
			},
			mocking: func(m *mocks.RealStateRepository, in input) output {
				m.
					On("GetRealState", mock.AnythingOfType("context.backgroundCtx"), in.id).
					Return(
						in.realState,
						nil,
					)

				m.
					On("UpdateRealState", mock.AnythingOfType("context.backgroundCtx"), in.realState, in.id).
					Return(
						domain.RealState{},
						errors.New("update real state failed"),
					)

				return output{
					realState: domain.RealState{},
					err:       errors.New("update real state failed"),
				}
			},
			assertion: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "When get fails, should return error",
			input: input{
				realState: domain.RealState{
					Registration: 987654321,
					Address:      "456 Elm St",
					Size:         200,
					Price:        275000.00,
					State:        "CA",
				},
				id: 1,
			},
			mocking: func(m *mocks.RealStateRepository, in input) output {
				m.
					On("GetRealState", mock.AnythingOfType("context.backgroundCtx"), in.id).
					Return(
						domain.RealState{},
						errors.New("get real state failed"),
					)

				return output{
					realState: domain.RealState{},
					err:       errors.New("get real state failed"),
				}
			},
			assertion: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected, actual)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run(tc.name, func(t *testing.T) {
				ctx := context.Background()

				r := mocks.NewRealStateRepository(t)
				s := service.NewRealStateService(r)

				expected := tc.mocking(r, tc.input)

				var actual output
				actual.realState, actual.err = s.Update(ctx, tc.input.realState, tc.input.id)

				tc.assertion(t, actual, expected)
			})
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		name      string
		input     uint64
		mocking   func(m *mocks.RealStateRepository, id uint64) error
		assertion func(t *testing.T, actual, expected error)
	}{
		{
			name:  "When real state exists, should delete and return nil",
			input: 1,
			mocking: func(m *mocks.RealStateRepository, id uint64) error {
				m.
					On("DeleteRealState", mock.AnythingOfType("context.backgroundCtx"), id).
					Return(nil)

				return nil
			},
			assertion: func(t *testing.T, actual, expected error) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:  "When real state exists, but delete fails should return error",
			input: 1,
			mocking: func(m *mocks.RealStateRepository, id uint64) error {
				m.
					On("DeleteRealState", mock.AnythingOfType("context.backgroundCtx"), id).
					Return(errors.New("delete real state failed"))

				return errors.New("delete real state failed")
			},
			assertion: func(t *testing.T, actual, expected error) {
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

			actual := s.Delete(ctx, tc.input)

			tc.assertion(t, actual, expected)
		})
	}
}
