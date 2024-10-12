package assertion

import (
	"github.com/stretchr/testify/assert"
)

// ErrorWithMessage can be used in table tests to check string equality of errors.
func ErrorWithMessage(msg string) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, resultErr error, msgAndArgs ...interface{}) bool {
		if resultErr == nil {
			return assert.Fail(t, "An error is expected, got nil.", msgAndArgs...)
		}

		return assert.EqualError(t, resultErr, msg, msgAndArgs...)
	}
}
