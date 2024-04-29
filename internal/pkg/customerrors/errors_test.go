package customerrors_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/natanchagas/gin-crud/internal/pkg/customerrors"
	"github.com/stretchr/testify/assert"
)

func TestCustomError(t *testing.T) {
	testCases := []struct {
		name      string
		err       error
		assertion func(t *testing.T, err error)
	}{
		{
			name: "When error is BadRequest, should be an customerror",
			err:  customerrors.BadRequest,
			assertion: func(t *testing.T, err error) {
				assert.IsType(t, customerrors.BadRequest, err)

				cerr, ok := err.(customerrors.Error)
				if !ok {
					t.Fail()
				}

				assert.Equal(t, http.StatusBadRequest, cerr.StatusCode)
				assert.Equal(t, customerrors.UserRequestError, cerr.ErrorCode)
			},
		},
		{
			name: "When error is not customerror, should be an error",
			err: fmt.Errorf("some error occurred"),
			assertion: func(t *testing.T, err error) {
                assert.IsType(t, fmt.Errorf(""), err)

				_, ok := err.(customerrors.Error)
				if ok {
					t.Fail()
				}
            },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertion(t, tc.err)
		})
	}
}
