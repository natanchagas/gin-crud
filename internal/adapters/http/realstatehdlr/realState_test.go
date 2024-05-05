package realstatehdlr_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
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
					httpCode: http.StatusCreated,
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
					httpCode: http.StatusInternalServerError,
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
					httpCode: http.StatusInternalServerError,
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
					httpCode: http.StatusBadRequest,
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

func TestGet(t *testing.T) {
	type output struct {
		httpCode int
		body     string
	}

	testCases := []struct {
		name       string
		input      string
		mocking    func(m *mocks.RealStateService, id string) output
		assertions func(t *testing.T, actual, expected output)
	}{
		{
			name:  "When input is valid and real state exists, should return real state with same id",
			input: "1",
			mocking: func(m *mocks.RealStateService, id string) output {
				rid, err := strconv.ParseUint(id, 10, 64)
				assert.NoError(t, err)

				realState := domain.RealState{
					Id:           rid,
					Registration: 987654321,
					Address:      "456 Elm St",
					Size:         200,
					Price:        250000.50,
					State:        "CA",
				}

				m.
					On("Get", mock.AnythingOfType("context.backgroundCtx"), rid).
					Return(realState, nil)

				b, err := json.Marshal(realState)
				assert.NoError(t, err)

				return output{
					httpCode: http.StatusOK,
					body:     string(b),
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:  "When input is invalid, should return 400",
			input: "a",
			mocking: func(m *mocks.RealStateService, id string) output {
				_, err := strconv.ParseUint(id, 10, 64)
				assert.Error(t, err)

				b, err := json.Marshal(customerrors.BadRequest)
				assert.NoError(t, err)

				return output{
					httpCode: http.StatusBadRequest,
					body:     string(b),
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:  "When input is valid and real state does not exists, should return 404",
			input: "1",
			mocking: func(m *mocks.RealStateService, id string) output {
				rid, err := strconv.ParseUint(id, 10, 64)
				assert.NoError(t, err)

				m.
					On("Get", mock.AnythingOfType("context.backgroundCtx"), rid).
					Return(domain.RealState{}, customerrors.NotFound)

				b, err := json.Marshal(customerrors.NotFound)
				assert.NoError(t, err)

				return output{
					httpCode: http.StatusNotFound,
					body:     string(b),
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:  "When input is valid and real state exists, but something goes wrong, should return 500",
			input: "1",
			mocking: func(m *mocks.RealStateService, id string) output {
				rid, err := strconv.ParseUint(id, 10, 64)
				assert.NoError(t, err)

				m.
					On("Get", mock.AnythingOfType("context.backgroundCtx"), rid).
					Return(domain.RealState{}, customerrors.Internal)

				b, err := json.Marshal(customerrors.Internal)
				assert.NoError(t, err)

				return output{
					httpCode: http.StatusInternalServerError,
					body:     string(b),
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:  "When input is valid and real state exists, but something unexpected goes wrong, should return 500",
			input: "1",
			mocking: func(m *mocks.RealStateService, id string) output {
				rid, err := strconv.ParseUint(id, 10, 64)
				assert.NoError(t, err)

				m.
					On("Get", mock.AnythingOfType("context.backgroundCtx"), rid).
					Return(domain.RealState{}, fmt.Errorf("unexpected error"))

				b, err := json.Marshal(customerrors.Unexpected)
				assert.NoError(t, err)

				return output{
					httpCode: http.StatusInternalServerError,
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
			req, _ := http.NewRequest("GET", fmt.Sprintf("/realstate/%s", tc.input), nil)
			router.ServeHTTP(w, req)

			var actual output
			actual.httpCode = w.Code
			actual.body = w.Body.String()

			tc.assertions(t, actual, expected)
		})
	}
}

func TestUpdate(t *testing.T) {
	type output struct {
		httpCode int
		body     string
	}

	type input struct {
		body string
		id   string
	}

	testCases := []struct {
		name       string
		input      input
		mocking    func(m *mocks.RealStateService, in input) output
		assertions func(t *testing.T, actual, expected output)
	}{
		{
			name: "When input is valid and real state exists, should return real state updated with same id",
			input: input{
				body: `{"registration": 987654321,"address": "456 Elm St","size": 200,"price": 275000.00,"state": "CA"}`,
				id:   "1",
			},
			mocking: func(m *mocks.RealStateService, in input) output {
				id, err := strconv.ParseUint(in.id, 10, 64)
				if err != nil {
					t.Fatal(err)
				}

				var realState domain.RealState

				err = json.Unmarshal([]byte(in.body), &realState)
				if err != nil {
					t.Fatal(err)
				}

				rs := realState
				rs.Id = id

				m.
					On("Update", mock.AnythingOfType("context.backgroundCtx"), realState, id).
					Return(rs, nil)

				b, err := json.Marshal(rs)
				if err != nil {
					t.Fatal(err)
				}

				return output{
					httpCode: http.StatusOK,
					body:     string(b),
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "When input is invalid, should return 400",
			input: input{
				body: `{"invalid"}`,
				id:   "1",
			},
			mocking: func(m *mocks.RealStateService, in input) output {
				_, err := strconv.ParseUint(in.id, 10, 64)
				if err != nil {
					t.Fatal(err)
				}

				b, err := json.Marshal(customerrors.BadRequest)
				if err != nil {
					t.Fatal(err)
				}

				return output{
					httpCode: http.StatusBadRequest,
					body:     string(b),
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "When input is valid and real state does not exists, should return 404",
			input: input{
				body: `{"registration": 987654321,"address": "456 Elm St","size": 200,"price": 275000.00,"state": "CA"}`,
				id:   "1",
			},
			mocking: func(m *mocks.RealStateService, in input) output {
				id, err := strconv.ParseUint(in.id, 10, 64)
				if err != nil {
					t.Fatal(err)
				}

				var realState domain.RealState

				err = json.Unmarshal([]byte(in.body), &realState)
				if err != nil {
					t.Fatal(err)
				}

				m.
					On("Update", mock.AnythingOfType("context.backgroundCtx"), realState, id).
					Return(domain.RealState{}, customerrors.NotFound)

				b, err := json.Marshal(customerrors.NotFound)
				if err != nil {
					t.Fatal(err)
				}

				return output{
					httpCode: http.StatusNotFound,
					body:     string(b),
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "When input is valid and real state exists, but something goes wrong, should return 500",
			input: input{
				body: `{"registration": 987654321,"address": "456 Elm St","size": 200,"price": 275000.00,"state": "CA"}`,
				id:   "1",
			},
			mocking: func(m *mocks.RealStateService, in input) output {
				id, err := strconv.ParseUint(in.id, 10, 64)
				if err != nil {
					t.Fatal(err)
				}

				var realState domain.RealState

				err = json.Unmarshal([]byte(in.body), &realState)
				if err != nil {
					t.Fatal(err)
				}

				m.
					On("Update", mock.AnythingOfType("context.backgroundCtx"), realState, id).
					Return(domain.RealState{}, customerrors.Internal)

				b, err := json.Marshal(customerrors.Internal)
				if err != nil {
					t.Fatal(err)
				}

				return output{
					httpCode: http.StatusInternalServerError,
					body:     string(b),
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "When input is valid and real state exists, but something unexpected goes wrong, should return 500",
			input: input{
				body: `{"registration": 987654321,"address": "456 Elm St","size": 200,"price": 275000.00,"state": "CA"}`,
				id:   "1",
			},
			mocking: func(m *mocks.RealStateService, in input) output {
				id, err := strconv.ParseUint(in.id, 10, 64)
				if err != nil {
					t.Fatal(err)
				}

				var realState domain.RealState

				err = json.Unmarshal([]byte(in.body), &realState)
				if err != nil {
					t.Fatal(err)
				}

				m.
					On("Update", mock.AnythingOfType("context.backgroundCtx"), realState, id).
					Return(domain.RealState{}, customerrors.Unexpected)

				b, err := json.Marshal(customerrors.Unexpected)
				if err != nil {
					t.Fatal(err)
				}

				return output{
					httpCode: http.StatusInternalServerError,
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
			gin.SetMode(gin.DebugMode)
			router := gin.Default()

			s := mocks.NewRealStateService(t)
			expected := tc.mocking(s, tc.input)

			hdl := realstatehdlr.NewRealStateHandler(s)

			hdl.BuildRoutes(router)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", fmt.Sprintf("/realstate/%s", tc.input.id), bytes.NewBuffer([]byte(tc.input.body)))
			router.ServeHTTP(w, req)

			var actual output
			actual.httpCode = w.Code
			actual.body = w.Body.String()

			tc.assertions(t, actual, expected)
		})
	}
}

func TestDelete(t *testing.T) {
	type output struct {
		httpCode int
		body     string
	}

	testCases := []struct {
		name       string
		input      string
		mocking    func(m *mocks.RealStateService, input string) output
		assertions func(t *testing.T, actual, expected output)
	}{
		{
			name:  "when input is valid and real state exists, should return 204",
			input: "1",
			mocking: func(m *mocks.RealStateService, input string) output {
				id, err := strconv.ParseUint(input, 10, 64)
				if err != nil {
					t.Fatal(err)
				}

				m.
					On("Delete", mock.AnythingOfType("context.backgroundCtx"), id).
					Return(nil)

				return output{
					httpCode: http.StatusNoContent,
					body:     "",
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:  "when input is invalid, should return 400",
			input: "a",
			mocking: func(m *mocks.RealStateService, input string) output {
				_, err := strconv.ParseUint(input, 10, 64)
				assert.Error(t, err)

				b, err := json.Marshal(customerrors.BadRequest)
				if err != nil {
					t.Fatal(err)
				}

				return output{
					httpCode: http.StatusBadRequest,
					body:     string(b),
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:  "when input is valid and real state exists, but something goes wrong, should return 500",
			input: "1",
			mocking: func(m *mocks.RealStateService, input string) output {
				id, err := strconv.ParseUint(input, 10, 64)
				if err != nil {
					t.Fatal(err)
				}

				m.
					On("Delete", mock.AnythingOfType("context.backgroundCtx"), id).
					Return(customerrors.Internal)

				b, err := json.Marshal(customerrors.Internal)
				if err != nil {
					t.Fatal(err)
				}

				return output{
					httpCode: http.StatusInternalServerError,
					body:     string(b),
				}
			},
			assertions: func(t *testing.T, actual, expected output) {
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:  "when input is valid and real state exists, but something goes unexpected wrong, should return 500",
			input: "1",
			mocking: func(m *mocks.RealStateService, input string) output {
				id, err := strconv.ParseUint(input, 10, 64)
				if err != nil {
					t.Fatal(err)
				}

				m.
					On("Delete", mock.AnythingOfType("context.backgroundCtx"), id).
					Return(fmt.Errorf("unexpected error"))

				b, err := json.Marshal(customerrors.Unexpected)
				if err != nil {
					t.Fatal(err)
				}

				return output{
					httpCode: http.StatusInternalServerError,
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
			gin.SetMode(gin.DebugMode)
			router := gin.Default()

			s := mocks.NewRealStateService(t)
			expected := tc.mocking(s, tc.input)

			hdl := realstatehdlr.NewRealStateHandler(s)

			hdl.BuildRoutes(router)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/realstate/%s", tc.input), nil)
			router.ServeHTTP(w, req)

			var actual output
			actual.httpCode = w.Code
			actual.body = w.Body.String()

			tc.assertions(t, actual, expected)
		})
	}

}
