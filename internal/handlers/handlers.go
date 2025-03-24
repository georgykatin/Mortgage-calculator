// Package handlers defines the HTTP request handlers for the mortgage calculation service.
//
// It contains logic for handling incoming requests, performing business logic for mortgage
// calculations, and interacting with the cache for storing and retrieving results.
// The main HTTP methods it handles are POST for performing mortgage calculations and GET for
// retrieving cached data.
//
// Functions and Methods:
//   - NewHandlers: Creates and returns a new Handlers instance with the provided cache storage.
//   - Execute: Handles the POST request for performing mortgage calculations.
//   - Cache: Handles the GET request for fetching cached data.
//   - monthlyPaymentCalculator: Calculates the monthly payment and overpayment for a mortgage.
//   - programValidator: Validates the selected loan program based on the input data.
//   - initialPaymentValidator: Validates that the initial payment is valid based on the object cost.
package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"sber/internal/cache"
	errs "sber/pkg/errors"
	"sber/pkg/models"
	"time"
)

// Handlers defines the HTTP request handlers for the mortgage calculation service.
// It stores a reference to the cache storage and provides methods to handle requests.
type Handlers struct {
	store *cache.Storage // The cache storage used for storing and retrieving mortgage calculation results
}

// NewHandlers creates a new Handlers instance with the provided cache storage.
func NewHandlers(store *cache.Storage) *Handlers {
	return &Handlers{store: store}
}

// Execute handles the POST request for performing mortgage calculations.
// It validates the input data, calculates the mortgage details, and stores the results in cache.
func (h *Handlers) Execute(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		err := json.NewEncoder(w).Encode(models.ErrorMessage{Error: "only post method allowed"})
		if err != nil {
			log.Println("failed to send error message in method check execute handler")
			return
		}
		return
	}

	// Decode the request body into ExecuteRequest structure
	reqData := models.ExecuteReqeust{}
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		log.Println("failed to decode request body in execute handler")
		return
	}

	// Validate the initial payment (should be at least 20% of the object cost)
	if !initialPaymentValidator(reqData.ObjectCost, reqData.InitialPayment) {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(models.ErrorMessage{Error: "the initial payment should be more"})
		if err != nil {
			log.Println("failed to send error message in method check execute handler")
			return
		}
		return
	}

	// Validate the selected loan program
	loanProgram, err := programValidator(reqData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		programValidatorErrorHandler(w, err)
		return
	}

	// Initialize rate and program based on the selected loan program
	rate, program := getLoanRateAndProgram(loanProgram)

	// Calculate monthly payment and overpayment
	monthlyPayment, overpayment := monthlyPaymentCalculator(float64(reqData.ObjectCost-reqData.InitialPayment), float64(rate), reqData.Months)

	// Prepare the response structure
	resp := prepareResponse(reqData, program, rate, monthlyPayment, overpayment)

	// Store the result in cache
	h.store.Load(resp.Result)

	// Send the response back to the client
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("failed to encode response message")
		return
	}
}

// Cache handles the GET request for fetching cached data.
// It retrieves data from the cache if available, otherwise, returns an error message.
func (h *Handlers) Cache(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		err := json.NewEncoder(w).Encode(models.ErrorMessage{Error: "only get method allowed"})
		if err != nil {
			log.Println("failed to send error message in method check execute handler")
			return
		}
		return
	}

	// Check if there is any data in the cache
	if !h.store.HasData() {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(models.ErrorMessage{Error: "empty cache"})
		if err != nil {
			log.Println("failed to send error message in cache data check ")
			return
		}
		return
	}

	// Retrieve all data from the cache
	data := h.store.ReadAll()

	// Send the cached data in the response
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}

// monthlyPaymentCalculator calculates the monthly payment and overpayment based on the loan amount,
// interest rate, and number of months for the mortgage.
func monthlyPaymentCalculator(objectCost, loanRate float64, months int32) (monthlyPayment, overpayment int32) {
	// Calculate the monthly interest rate
	monthlyRate := loanRate / (100 * 12)

	// Calculate the factor for the loan formula
	factor := math.Pow((1 + monthlyRate), float64(months))

	// Calculate the monthly payment
	monthlyPayment = int32(math.Ceil(objectCost * (monthlyRate * factor) / (factor - 1)))

	// Calculate the overpayment
	overpayment = int32(math.Ceil(float64(monthlyPayment)*float64(months) - objectCost))

	// Return the calculated values as integers
	return monthlyPayment, overpayment
}

// programValidator validates the loan program based on the request data.
// It checks if exactly one program flag is set to true, and returns an error if any validation fails.
func programValidator(data models.ExecuteReqeust) (string, error) {
	// Declare variables for error handling and to hold the result of validation checks.
	var (
		countTrue     int    // Tracks how many program flags are set to true
		lastTrueField string // Holds the name of the last program flag that was set to true
		hasAnyField   bool   // Indicates whether at least one program flag is set to true
	)

	// Check each program flag manually
	if data.Program.Base {
		countTrue++
		lastTrueField = "base"
	}
	if data.Program.Military {
		countTrue++
		lastTrueField = "military"
	}
	if data.Program.Salary {
		countTrue++
		lastTrueField = "salary"
	}

	// Check if at least one program flag is set to true
	if data.Program.Base || data.Program.Military || data.Program.Salary {
		hasAnyField = true
	}

	// Return errors if no flags or more than one flag are set to true
	if !hasAnyField {
		return "", errs.ErrNoTrueValues
	}
	if countTrue == 0 {
		return "", errs.ErrNoTrueValues
	}
	if countTrue > 1 {
		return "", errs.ErrMoreThanOneTrue
	}

	// Return the name of the program flag that was set to true
	return lastTrueField, nil
}

// InitialPaymentValidator validates the initial payment based on the object cost.
// The initial payment must be more than zero and at least 20% of the object cost.
func initialPaymentValidator(objectCost, initialPayment int32) bool {
	if initialPayment > objectCost {
		return false
	}
	// If both object cost and initial payment are zero, return false
	if objectCost == 0 && initialPayment == 0 {
		return false
	}
	// If the initial payment is zero, return false
	if initialPayment == 0 {
		return false
	}
	// If the initial payment is less than 20% of the object cost, return false
	if initialPayment*5 < objectCost {
		return false
	}

	// If all conditions are satisfied, return true
	return true
}

func programValidatorErrorHandler(w http.ResponseWriter, err error) {
	if errors.Is(err, errs.ErrNoTrueValues) {
		err = json.NewEncoder(w).Encode(models.ErrorMessage{Error: "choose program"})
		if err != nil {
			log.Println("failed to send error message in program validator error handler")
			return
		}
		return
	}
	if errors.Is(err, errs.ErrMoreThanOneTrue) {
		err = json.NewEncoder(w).Encode(models.ErrorMessage{Error: "choose only 1 program"})
		if err != nil {
			log.Println("failed to send error message in program validator error handler")
			return
		}
		return
	}
}

func prepareResponse(reqData models.ExecuteReqeust, program models.Program, rate uint8, monthlyPayment, overpayment int32) models.ExecuteResponse {
	lastDate := time.Now().AddDate(0, int(reqData.Months), 0).Format("2006-01-01")
	return models.ExecuteResponse{
		Result: models.Result{
			Params: models.Params{
				ObjectCost:     reqData.ObjectCost,
				InitialPayment: reqData.InitialPayment,
				Months:         reqData.Months,
			},
			Program: program,
			Aggregates: models.Aggregates{
				Rate:            rate,
				LoanSum:         reqData.ObjectCost - reqData.InitialPayment,
				MonthlyPayment:  monthlyPayment,
				Overpayment:     overpayment,
				LastPaymentDate: lastDate,
			},
		},
	}
}

func getLoanRateAndProgram(loanProgram string) (uint8, models.Program) {
	var rate uint8
	program := models.Program{}
	switch loanProgram {
	case "base":
		rate = 10
		program = models.Program{Base: true}
	case "military":
		rate = 9
		program = models.Program{Military: true}
	case "salary":
		rate = 8
		program = models.Program{Salary: true}
	}
	return rate, program
}
