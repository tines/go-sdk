package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tines/go-sdk/internal/utils"
)

func TestSetUserAgent(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		input       string
		expected    string
		expectedErr error
	}{
		{input: "", expected: "Tines/GoSdk", expectedErr: nil},
		{input: "foo", expected: "foo", expectedErr: nil},
	}

	for _, tt := range tests {
		actual := utils.SetUserAgent(tt.input)
		assert.Equal(tt.expected, actual)
	}
}
