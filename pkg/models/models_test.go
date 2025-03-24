package models

import (
	"encoding/json"
	"testing"
)

func TestCacheStorageFormat_MarshalJSON(t *testing.T) {
	// Input data
	data := CacheStorageFormat{
		ID: 1,
		Params: Params{
			ObjectCost:     5000000,
			InitialPayment: 1000000,
			Months:         240,
		},
		Program: Program{
			Salary: true,
		},
		Aggregates: Aggregates{
			Rate:            8,
			LoanSum:         4000000,
			MonthlyPayment:  33458,
			Overpayment:     4029920,
			LastPaymentDate: "2044-02-18",
		},
	}

	// Expected JSON output (strict field order!)
	expectedJSON := `{
		"id": 1,
		"params": {
			"object_cost": 5000000,
			"initial_payment": 1000000,
			"months": 240
		},
		"program": {
			"salary": true
		},
		"aggregates": {
			"rate": 8,
			"loan_sum": 4000000,
			"monthly_payment": 33458,
			"overpayment": 4029920,
			"last_payment_date": "2044-02-18"
		}
	}`

	// Serialize the struct
	resultJSON, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Remove extra spaces and compare
	var expected, result map[string]interface{}
	if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}
	if err := json.Unmarshal(resultJSON, &result); err != nil {
		t.Fatalf("Failed to unmarshal result JSON: %v", err)
	}

	if len(expected) != len(result) {
		t.Errorf("JSON field count mismatch: expected %d, got %d", len(expected), len(result))
	}

	// Check if JSON outputs match
	if !jsonEqual(expected, result) {
		t.Errorf("Mismatch in JSON output.\nExpected: %s\nGot: %s", expectedJSON, string(resultJSON))
	}
}

// jsonEqual compares two JSON objects (maps).
func jsonEqual(a, b map[string]interface{}) bool {
	aBytes, _ := json.Marshal(a)
	bBytes, _ := json.Marshal(b)
	return string(aBytes) == string(bBytes)
}
