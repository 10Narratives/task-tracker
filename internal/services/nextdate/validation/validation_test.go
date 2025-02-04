package validation_test

import (
	"testing"

	"github.com/10Narratives/task-tracker/internal/services/nextdate/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		timeStep string
		pattern  string
		wantErr  require.ErrorAssertionFunc
	}{
		{"Valid timeStep for days", "d 7", "^d (?:[1-9]|[1-9][0-9]|[1-3][0-9]{2}|400)$", require.NoError},
		{"Invalid timeStep for days", "d 700", "^d (?:[1-9]|[1-9][0-9]|[1-3][0-9]{2}|400)$", func(tt require.TestingT, err error, i ...interface{}) {
			assert.EqualError(t, err, "invalid time step format")
		}},
		{"Valid timeStep for years", "y", "^y$", require.NoError},
		{"Invalid timeStep for years", "y 100", "^y$", func(tt require.TestingT, err error, i ...interface{}) {
			assert.EqualError(t, err, "invalid time step format")
		}},
		{"Valid timeStep for moths", "m 7", "^m (-?[1-9]|-1|-2|[12][0-9]|3[01])(,(-?[1-9]|-1|-2|[12][0-9]|3[01]))*( (1[0-2]|[1-9])(,(1[0-2]|[1-9]))*)?$", require.NoError},
		{"Invalid timeStep for moths", "m 700", "^m (-?[1-9]|-1|-2|[12][0-9]|3[01])(,(-?[1-9]|-1|-2|[12][0-9]|3[01]))*( (1[0-2]|[1-9])(,(1[0-2]|[1-9]))*)?$", func(tt require.TestingT, err error, i ...interface{}) {
			assert.EqualError(t, err, "invalid time step format")
		}},
		{"Valid timeStep for weeks", "w 7", "^w ([1-7](,[1-7])*)$", require.NoError},
		{"Invalid timeStep for weeks", "w 700", "^w ([1-7](,[1-7])*)$", func(tt require.TestingT, err error, i ...interface{}) {
			assert.EqualError(t, err, "invalid time step format")
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.Validate(tt.timeStep, tt.pattern)
			tt.wantErr(t, err)
		})
	}
}
