package hw09structvalidator

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version  string `validate:"len:5"`
		Response Response
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
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "successful case for user",
			in: User{
				ID:     "2bb933f2-8de1-4172-97e8-2de324301317",
				Name:   "Test",
				Age:    20,
				Email:  "test@test.ru",
				Role:   "admin",
				Phones: []string{"12345678999"},
			},
			expectedErr: nil,
		},
		{
			name: "failed with not existed role",
			in: User{
				Role: "notExist",
			},
			expectedErr: ErrorNotInSlice,
		},
		{
			name: "failed with value more than max",
			in: User{
				Age: 55,
			},
			expectedErr: ErrorMoreThanMax,
		},
		{
			name: "failed with value less than min",
			in: User{
				Age: 5,
			},
			expectedErr: ErrorLessThanMin,
		},
		{
			name: "failed with does not match len",
			in: User{
				Phones: []string{"1234567899922222222222222"},
			},
			expectedErr: ErrorLessThanMin,
		},
		{
			name: "failed with value does not match regexp",
			in: User{
				Email: "test@test.ru.ru",
			},
			expectedErr: ErrorDoesNotMatchRegularExpression,
		},
		{
			name: "failed with value is not in slice in a nested structure",
			in: App{
				Version: "1.005",
				Response: Response{
					Code: 401,
				},
			},
			expectedErr: ErrorNotInSlice,
		},
		{
			name: "successful case for structure without validation tags",
			in: Token{
				Header: []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"),
				Payload: []byte("eyJhdWQiOiJyZWZyZXNoIiwiZXhwIjoxNzM0NDI0OTQ5LCJoYXNoIjoi" +
					"YmVlMDRjZGEtNDllNS00OTA1LTg0YTItZDdkZTBjN2RlMDY4Iiwicm9sZSI6IlNVUEVSX0FETU" +
					"lOIiwic3RhdGUiOiJBQ1RJVkUiLCJzdWIiOiIyODY3YTE5MS1lMDRlLTQyY2QtOWM1YS1jMzM3MGQ5OTE0NmQifQ"),
				Signature: []byte("QMenFBGoWgYBNT7bdyaAUodQ957b_AcHDAc6gLltJkU"),
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := ValidateStruct(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.expectedErr.Error())
			}
		})
	}
}
