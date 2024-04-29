package realstatehdlr_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/natanchagas/gin-crud/internal/adapters/http/realstatehdlr"
	"github.com/natanchagas/gin-crud/internal/core/domain"
	"github.com/natanchagas/gin-crud/internal/mocks"
	"github.com/natanchagas/gin-crud/internal/pkg/customerrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	type output struct {
		httpCode int
		body     string
	}

	testCases := []struct {
		name       string
		input      string
		mocking    func(m *mocks.RealStateService, realState string) output
		assertions func(t *testing.T, actual, expected output)
	}{
		{
			name:  "When payload is valid and service returns success, should return real state with id",
			input: `{"registration":987654321,"address":"456 Elm St","size":200,"price":250000.5,"state":"CA"}`,
			mocking: func(m *mocks.RealStateService, realState string) output {
				var rs domain.RealState
				var ors domain.RealState

				err := json.Unmarshal([]byte(realState), &rs)
				if err != nil {
					t.Fatal(err)
				}

				ors = rs
				ors.Id = 1

				m.
					On("Create", mock.AnythingOfType("context.backgroundCtx"), rs).
					Return(ors, nil)

				b, err := json.Marshal(ors)
				if err != nil {
					t.Fatal(err)
				}

				return output{
					httpCode: 201,
					body:     string(b),
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected.httpCode, actual.httpCode)
				assert.Equal(t, expected.body, actual.body)
			},
		},
		{
			name:  "When payload is valid, but service fails, should return error",
			input: `{"registration":987654321,"address":"456 Elm St","size":200,"price":250000.5,"state":"CA"}`,
			mocking: func(m *mocks.RealStateService, realState string) output {
				var rs domain.RealState

				err := json.Unmarshal([]byte(realState), &rs)
				if err != nil {
					t.Fatal(err)
				}

				m.
					On("Create", mock.AnythingOfType("context.backgroundCtx"), rs).
					Return(domain.RealState{}, customerrors.Internal)

				b, err := json.Marshal(customerrors.Internal)
				assert.NoError(t, err)

				return output{
					httpCode: 500,
					body:     string(b),
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:  "When payload is valid, but service fails with an unexpected error, should return error",
			input: `{"registration":987654321,"address":"456 Elm St","size":200,"price":250000.5,"state":"CA"}`,
			mocking: func(m *mocks.RealStateService, realState string) output {
				var rs domain.RealState

				err := json.Unmarshal([]byte(realState), &rs)
				if err != nil {
					t.Fatal(err)
				}

				m.
					On("Create", mock.AnythingOfType("context.backgroundCtx"), rs).
					Return(domain.RealState{}, errors.New("unexpected error"))

				b, err := json.Marshal(customerrors.Unexpected)
				assert.NoError(t, err)

				return output{
					httpCode: 500,
					body:     string(b),
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:  "When payload is invalid, should return error",
			input: `{"invalid"}`,
			mocking: func(m *mocks.RealStateService, realState string) output {

				b, err := json.Marshal(customerrors.BadRequest)
				assert.NoError(t, err)

				return output{
					httpCode: 400,
					body:     string(b),
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected, actual)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			gin.SetMode(gin.TestMode)
			router := gin.Default()

			s := mocks.NewRealStateService(t)

			expected := tc.mocking(s, tc.input)

			hdlr := realstatehdlr.NewRealStateHandler(s)
			hdlr.BuildRoutes(router)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/realstate/", bytes.NewBuffer([]byte(tc.input)))
			router.ServeHTTP(w, req)

			var actual output
			actual.httpCode = w.Code
			actual.body = w.Body.String()

			tc.assertions(t, actual, expected)
		})
	}
}
