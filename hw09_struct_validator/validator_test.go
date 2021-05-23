package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in              interface{}
		expectedErr     error
		expectedValErrs []error
	}{
		{
			User{"100", "testname", 10, "test@test.ru", "stuff", []string{"5466", "6546458483"}, nil},
			ErrInvalidValues,
			[]error{ErrLen, ErrMin, ErrLen, ErrLen},
		},
		{App{"112783"}, ErrInvalidValues, []error{ErrLen}},
		{Token{}, nil, nil},
		{Response{200, "test"}, nil, nil},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tc := tc
			t.Parallel()

			valErrs, err := Validate(tc.in)

			require.ErrorIs(t, tc.expectedErr, err, "Error should be like expected")

			for i, err := range tc.expectedValErrs {
				require.ErrorIs(t, valErrs[i].Err, err, "Validation error should be like expected")
			}
		})
	}

	t.Run("should handle non struct value", func(t *testing.T) {
		_, err := Validate(123)

		require.ErrorIs(t, err, ErrExpectedStruct, "Should throw nonStruct error")
	})

	t.Run("should handle nil value", func(t *testing.T) {
		valErrs, err := Validate(nil)

		require.ErrorIs(t, err, ErrExpectedStruct, "Should throw nonStruct error")
		require.Nil(t, valErrs, "Val Errors should be nil")
	})
}
