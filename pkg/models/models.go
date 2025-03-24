// Package models defines the data structures used for mortgage calculations,
// request/response handling, and caching in the application. The package includes
// various structs representing parameters, programs, results, and aggregated information
// related to mortgage calculations, as well as structures for handling API responses
// and errors.
//
// Calculated parameters include:
//   - Interest rate based on the requested loan program
//   - Loan amount
//   - Annuity monthly payment
//   - Total overpayment over the loan period
//   - Last payment date
package models

import "encoding/json"

// Aggregates represents the calculated financial aggregates based on the mortgage request.
// It includes the interest rate, loan sum, monthly payment, total overpayment, and the last payment date.
type Aggregates struct {
	LastPaymentDate string `json:"last_payment_date"` // Date of the last payment
	Rate            uint8  `json:"rate"`              // Interest rate
	LoanSum         int32  `json:"loan_sum"`          // Loan amount
	MonthlyPayment  int32  `json:"monthly_payment"`   // Monthly payment amount
	Overpayment     int32  `json:"overpayment"`       // Total overpayment for the loan
}

// Program represents different mortgage programs with flags indicating whether they
// apply to salary-based, military, or base programs.
type Program struct {
	Salary   bool `json:"salary,omitempty"`   // Indicates if the program is salary-based
	Military bool `json:"military,omitempty"` // Indicates if the program is military
	Base     bool `json:"base,omitempty"`     // Indicates if the program is base-based
}

// Params contains the core parameters needed for mortgage calculations such as
// object cost, initial payment, and the loan term in months.
type Params struct {
	ObjectCost     int32 `json:"object_cost"`     // The cost of the object being purchased
	InitialPayment int32 `json:"initial_payment"` // The initial payment amount
	Months         int32 `json:"months"`          // Loan term in months
}

// ExecuteReqeust represents the structure of a request to execute the mortgage calculation.
// It contains the object cost, initial payment, loan term, and program details.
type ExecuteReqeust struct {
	ObjectCost     int32   `json:"object_cost"`     // Object cost for the loan
	InitialPayment int32   `json:"initial_payment"` // Initial payment amount
	Months         int32   `json:"months"`          // Loan term in months
	Program        Program `json:"program"`         // Mortgage program details
}

// ExecuteResponse represents the structure of the response containing the mortgage calculation result.
type ExecuteResponse struct {
	Result Result `json:"result"` // The result of the mortgage calculation
}

// Result contains the detailed mortgage calculation results, including parameters, the program,
// and the aggregated financial data (interest rate, loan sum, etc.).
type Result struct {
	Aggregates Aggregates `json:"aggregates"` // Calculated aggregates (interest rate, overpayment, etc.)
	Params     Params     `json:"params"`     // Mortgage parameters (object cost, initial payment, etc.)
	Program    Program    `json:"program"`    // Mortgage program (salary, military, base, etc.)

}

// CacheStorageFormat represents the structure of a cached mortgage calculation.
// It stores the ID, parameters, program, and calculated aggregates.
type CacheStorageFormat struct {
	Aggregates Aggregates `json:"aggregates"` // Calculated aggregates (interest rate, overpayment, etc.)
	Params     Params     `json:"params"`     // Mortgage parameters
	Program    Program    `json:"program"`    // Mortgage program details
	ID         int32      `json:"id"`         // Unique identifier for the cached entry
}

// CacheResponse is the structure for returning a list of cached mortgage calculations.
type CacheResponse struct {
	Results []CacheStorageFormat // List of cached mortgage calculations
}

// ErrorMessage represents an error message returned by the API.
type ErrorMessage struct {
	Error string `json:"error"` // The error message
}

// The MarshalJSON methods preserves the alignment of fields in the underlying CacheStorageFormat structure,
//
//	but at the same time assembles the JSON in the order required for storage
func (c CacheStorageFormat) MarshalJSON() ([]byte, error) {
	type Alias CacheStorageFormat
	return json.Marshal(&struct {
		ID         int32      `json:"id"`
		Params     Params     `json:"params"`
		Program    Program    `json:"program"`
		Aggregates Aggregates `json:"aggregates"`
		*Alias
	}{
		ID:         c.ID,
		Params:     c.Params,
		Program:    c.Program,
		Aggregates: c.Aggregates,
		Alias:      (*Alias)(&c),
	})
}
