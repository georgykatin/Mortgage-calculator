// Package errors defines custom error types used throughout the application.
// These errors are specifically used for validation in different parts of the loan program
// and initial payment processing. The errors help provide clear, meaningful messages
// for various validation failures.
package errors

import "errors"

// Custom errors for loan program validation.
var (
	// ErrNoTrueValues is returned when no true values are found in a set of values
	// that are expected to contain at least one true value.
	ErrNoTrueValues = errors.New("there is no true values")

	// ErrMoreThanOneTrue is returned when more than one true value is found in a set
	// where only one true value is expected.
	ErrMoreThanOneTrue = errors.New("there are more that one true value")
)

// Custom errors for initial payment validation.
var (
	// ErrInitalPaymentIsTooSmall is returned when the initial payment is too small
	// and doesn't meet the minimum required amount.
	ErrInitalPaymentIsTooSmall = errors.New("the initial payment should be more")
)

// Custom errors for cofig load.
var (
	// ErrInvalidPath is returned when file with filepath is not in safe directory.
	ErrInvalidPath = errors.New("invalid path")
)
