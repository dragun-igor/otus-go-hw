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
		ID     string `json:"id" validate:"len:10"`
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
		Header    []byte `validate:"len:1"`
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	WrongTag struct {
		ID string `validate:"WRONGTAG:20"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in                    interface{}
		expectedValidationErr ValidationErrors
		expectedErr           error
	}{
		// case 0
		{
			in:                    []string{"0", "1", "2"},
			expectedValidationErr: ValidationErrors{},
			expectedErr:           ErrNotStruct,
		},
		// case 1
		{
			in:                    App{"1234"},
			expectedValidationErr: ValidationErrors{ValidationError{"Version", ErrLen}},
			expectedErr:           nil,
		},
		// case 2
		{
			in:                    App{"12345"},
			expectedValidationErr: ValidationErrors{},
			expectedErr:           nil,
		},
		// case 3
		{
			in: User{
				ID:     "12345",
				Name:   "Igor",
				Age:    17,
				Email:  "igor@mail",
				Role:   "guest",
				Phones: []string{"9875556677", "89123456789", "9123456789"},
				meta:   []byte{1, 2, 3},
			},
			expectedValidationErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrLen,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrMin,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrRegexp,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrIn,
				},
				ValidationError{
					Field: "Phones[0]",
					Err:   ErrLen,
				},
				ValidationError{
					Field: "Phones[2]",
					Err:   ErrLen,
				},
			},
			expectedErr: nil,
		},
		// case 4
		{
			in: Token{
				Header: []byte{1, 2, 3},
			},
			expectedValidationErr: ValidationErrors{},
			expectedErr:           ErrUnknownType,
		},
		// case 5
		{
			in: Response{
				Code: 900,
				Body: "1234",
			},
			expectedValidationErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrIn,
				},
			},
			expectedErr: nil,
		},
		// case 6
		{
			in: WrongTag{
				ID: "123",
			},
			expectedValidationErr: ValidationErrors{},
			expectedErr:           ErrUnknownTag,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			valErr, err := Validate(ValidationErrors{}, tt.in)
			require.Equal(t, tt.expectedValidationErr, valErr)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
