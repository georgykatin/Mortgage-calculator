package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sber/internal/cache"
	"sber/pkg/models"
	"testing"
	"time"
)

func TestExecuteHandler(t *testing.T) {
	mockCache := cache.New()
	h := NewHandlers(mockCache)

	tests := []struct {
		name         string
		payload      any
		Method       string
		expectedCode int
		checkCache   bool
	}{
		{"Valid request",
			models.ExecuteReqeust{
				ObjectCost:     100000,
				InitialPayment: 20000,
				Months:         12,
				Program:        models.Program{Base: true},
			},
			"POST",
			http.StatusOK,
			true,
		},
		{"Invalid method",
			"",
			"GET",
			http.StatusMethodNotAllowed,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(tt.Method, "/execute", bytes.NewReader(body))
			w := httptest.NewRecorder()

			h.Execute(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.checkCache && !mockCache.HasData() {
				t.Error("Expected data in cache but found none")
			}
		})
	}
}

func TestCacheHandler(t *testing.T) {
	mockCache := cache.New()
	h := NewHandlers(mockCache)

	// Prepopulate cache
	mockCache.Load(models.Result{
		Params: models.Params{
			ObjectCost:     100000,
			InitialPayment: 20000,
			Months:         12,
		},
		Aggregates: models.Aggregates{
			Rate:            10,
			LoanSum:         80000,
			MonthlyPayment:  8792,
			Overpayment:     5504,
			LastPaymentDate: time.Now().AddDate(0, 12, 0).Format("2006-01-02"),
		},
	})

	t.Run("Valid cache request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/cache", nil)
		w := httptest.NewRecorder()

		h.Cache(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response []models.Result
		json.NewDecoder(w.Body).Decode(&response)

		if len(response) != 1 {
			t.Errorf("Expected 1 result, got %d", len(response))
		}
	})
}

func TestCacheHandlerEmpty(t *testing.T) {
	// Создаем чистый кеш
	mockCache := cache.New()
	h := NewHandlers(mockCache)

	t.Run("Empty cache request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/cache", nil)
		w := httptest.NewRecorder()

		h.Cache(w, req)

		// Проверяем статус код
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}

		// Проверяем тело ответа
		var errMsg models.ErrorMessage
		json.NewDecoder(w.Body).Decode(&errMsg)

		expectedError := "empty cache"
		if errMsg.Error != expectedError {
			t.Errorf("Expected error '%s', got '%s'", expectedError, errMsg.Error)
		}
	})
}
