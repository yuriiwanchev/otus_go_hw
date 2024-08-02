package hw09structvalidator

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
		Meta   Meta     `validate:"nested"`
	}

	Meta struct {
		Description string `validate:"len:1"`
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
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:    "12345678-1234-1234-1234-123456789012",
				Name:  "John",
				Age:   20,
				Email: "test@example.com",
				Role:  "admin",
				Phones: []string{
					"12345678901",
				},
				Meta: Meta{
					Description: "t",
				},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:    "short-id",
				Name:  "John",
				Age:   17,
				Email: "invalid-email",
				Role:  "user",
				Phones: []string{
					"short",
				},
				Meta: Meta{
					Description: "test",
				},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: fmt.Errorf("must be 36 characters long")},
				{Field: "Age", Err: fmt.Errorf("must be at least 18")},
				{Field: "Email", Err: fmt.Errorf("must match regexp ^\\w+@\\w+\\.\\w+$")},
				{Field: "Role", Err: fmt.Errorf("must be one of [admin stuff]")},
				{Field: "Phones", Err: fmt.Errorf("element 0: must be 11 characters long")},
				{
					Field: "Meta", Err: ValidationErrors{
						{Field: "Description", Err: fmt.Errorf("must be 1 characters long")},
					},
				},
			},
		},
		{
			in: User{
				ID:    "12345678-1234-1234-1234-123456789012",
				Name:  "John",
				Age:   55,
				Email: "test@example.com",
				Role:  "admin",
				Phones: []string{
					"12345678901",
				},
				Meta: Meta{
					Description: "t",
				},
			},
			expectedErr: ValidationErrors{
				{Field: "Age", Err: fmt.Errorf("must be at most 50")},
			},
		},
		{
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1.0.0.0",
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: fmt.Errorf("must be 5 characters long")},
			},
		},
		{
			in: Response{
				Code: 200,
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 201,
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: fmt.Errorf("must be one of [200 404 500]")},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)

			if tt.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.IsType(t, tt.expectedErr, err)

				var expectedValidationErrors *ValidationErrors
				_ = errors.As(tt.expectedErr, &expectedValidationErrors)

				var actualValidationErrors *ValidationErrors
				ok := errors.As(err, &actualValidationErrors)

				if ok {
					assert.Equal(t, len(*expectedValidationErrors), len(*actualValidationErrors))

					for _, expectedErr := range *expectedValidationErrors {
						found := false
						for _, actualErr := range *actualValidationErrors {
							if actualErr.Field == expectedErr.Field && actualErr.Err.Error() == expectedErr.Err.Error() {
								found = true
								break
							}
						}
						assert.True(t, found, fmt.Sprintf("expected error for field %s not found", expectedErr.Field))
					}
				}
			}
		})
	}
}
