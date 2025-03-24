package handlers

import (
	"errors"
	"reflect"
	errs "sber/pkg/errors"
	"sber/pkg/models"
	"testing"
)

func TestInitialPaymentValidator(t *testing.T) {
	tests := []struct {
		name        string
		objectCost  int32
		initialPay  int32
		expectValid bool
	}{
		{"Zero values", 0, 0, false},
		{"Initial pay zero", 100000, 0, false},
		{"Valid 20% payment", 100000, 20000, true},
		{"Payment below 20%", 100000, 19000, false},
		{"Payment exceeds cost", 100000, 150000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res := initialPaymentValidator(tt.objectCost, tt.initialPay); res != tt.expectValid {
				t.Errorf("Expected %v, got %v", tt.expectValid, res)
			}
		})
	}
}

func TestProgramValidator(t *testing.T) {
	tests := []struct {
		name        string
		program     models.Program
		expectError error
		expectName  string
	}{
		{"No program selected", models.Program{}, errs.ErrNoTrueValues, ""},
		{"Multiple programs",
			models.Program{Base: true, Military: true},
			errs.ErrMoreThanOneTrue, ""},
		{"No true programs", models.Program{Base: false, Military: false},
			errs.ErrNoTrueValues, ""},
		{"Valid base program",
			models.Program{Base: true}, nil, "base"},
		{"Valid military program",
			models.Program{Military: true}, nil, "military"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := models.ExecuteReqeust{Program: tt.program}
			name, err := programValidator(req)

			if !errors.Is(err, tt.expectError) {
				t.Errorf("Expected error %v, got %v", tt.expectError, err)
			}
			if name != tt.expectName {
				t.Errorf("Expected name %s, got %s", tt.expectName, name)
			}
		})
	}
}

func TestMonthlyPaymentCalculator(t *testing.T) {
	tests := []struct {
		name            string
		objectCost      float64
		loanRate        float64
		months          int32
		expectedPayment int32
		expectedOverpay int32
	}{
		{
			name:            "Basic calculation",
			objectCost:      100000,
			loanRate:        10,
			months:          12,
			expectedPayment: 8792,
			expectedOverpay: 5504,
		},
		{
			name:            "Short term loan",
			objectCost:      50000,
			loanRate:        5,
			months:          6,
			expectedPayment: 8456,
			expectedOverpay: 736,
		},
		{
			name:            "Long term loan",
			objectCost:      200000,
			loanRate:        7.5,
			months:          240,
			expectedPayment: 1612,
			expectedOverpay: 186880,
		},
		{
			name:            "Small loan amount",
			objectCost:      1000,
			loanRate:        5,
			months:          12,
			expectedPayment: 86,
			expectedOverpay: 32,
		},
		{
			name:            "High interest rate",
			objectCost:      100000,
			loanRate:        20,
			months:          12,
			expectedPayment: 9264,
			expectedOverpay: 11168,
		},
		{
			name:            "One month term",
			objectCost:      10000,
			loanRate:        10,
			months:          1,
			expectedPayment: 10084,
			expectedOverpay: 84,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment, overpay := monthlyPaymentCalculator(tt.objectCost, tt.loanRate, tt.months)

			if payment != tt.expectedPayment {
				t.Errorf("Expected monthly payment %d, got %d", tt.expectedPayment, payment)
			}
			if overpay != tt.expectedOverpay {
				t.Errorf("Expected overpayment %d, got %d", tt.expectedOverpay, overpay)
			}
		})
	}
}

func TestGetLoanRateAndProgram(t *testing.T) {
	tests := []struct {
		name            string
		loanProgram     string
		expectedRate    uint8
		expectedProgram models.Program
	}{
		{
			name:            "Base program",
			loanProgram:     "base",
			expectedRate:    10,
			expectedProgram: models.Program{Base: true},
		},
		{
			name:            "Military program",
			loanProgram:     "military",
			expectedRate:    9,
			expectedProgram: models.Program{Military: true},
		},
		{
			name:            "Salary program",
			loanProgram:     "salary",
			expectedRate:    8,
			expectedProgram: models.Program{Salary: true},
		},
		{
			name:            "Unknown program",
			loanProgram:     "unknown",
			expectedRate:    0,
			expectedProgram: models.Program{}, // Пустая структура, так как программа неизвестна
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rate, program := getLoanRateAndProgram(tt.loanProgram)

			if rate != tt.expectedRate {
				t.Errorf("got rate %d, want %d", rate, tt.expectedRate)
			}
			if !reflect.DeepEqual(program, tt.expectedProgram) {
				t.Errorf("got program %+v, want %+v", program, tt.expectedProgram)
			}
		})
	}
}
