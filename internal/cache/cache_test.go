package cache_test

import (
	"sber/internal/cache"
	"sber/pkg/models"
	"sync"
	"sync/atomic"
	"testing"
)

// TestNew verifies that the New function creates a new Storage instance with an empty cache.
func TestNew(t *testing.T) {
	storage := cache.New()
	if storage == nil {
		t.Error("Expected new Storage instance, got nil")
	}
	if len(storage.ReadAll()) != 0 {
		t.Error("Expected empty cache, got non-empty cache")
	}
	if storage.HasData() {
		t.Error("Expected cache to be empty, but HasData() returned true")
	}
}

// TestLoad verifies that the Load method adds data to the cache and increments IDCounter correctly.
func TestLoad(t *testing.T) {
	storage := cache.New()

	// Creating test data
	result := models.Result{
		Params: models.Params{
			ObjectCost:     100000,
			InitialPayment: 20000,
			Months:         12,
		},
		Program: models.Program{Base: true},
		Aggregates: models.Aggregates{
			Rate:            10,
			LoanSum:         80000,
			MonthlyPayment:  8792,
			Overpayment:     5504,
			LastPaymentDate: "2024-01-01",
		},
	}

	// Adding data to cache
	storage.Load(result)

	// Verifying that data has been added
	cachedData := storage.ReadAll()
	if len(cachedData) != 1 {
		t.Errorf("Expected 1 entry in cache, got %d", len(cachedData))
	}
	if cachedData[0].ID != 0 {
		t.Errorf("Expected ID to be 0, got %d", cachedData[0].ID)
	}
	if cachedData[0].Params != result.Params {
		t.Error("Cached Params do not match input")
	}
	if cachedData[0].Program != result.Program {
		t.Error("Cached Program does not match input")
	}
	if cachedData[0].Aggregates != result.Aggregates {
		t.Error("Cached Aggregates do not match input")
	}

	// Verifying that IDCounter has incremented
	storage.Load(result)
	cachedData = storage.ReadAll()
	if len(cachedData) != 2 {
		t.Errorf("Expected 2 entries in cache, got %d", len(cachedData))
	}
	if cachedData[1].ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", cachedData[1].ID)
	}
}

// TestReadAll verifies that the ReadAll method returns all data from the cache.
func TestReadAll(t *testing.T) {
	storage := cache.New()

	// Creating test data
	testData := []models.Result{
		{Params: models.Params{ObjectCost: 100000}},
		{Params: models.Params{ObjectCost: 200000}},
		{Params: models.Params{ObjectCost: 300000}},
	}

	// Adding data to cache
	for _, data := range testData {
		storage.Load(data)
	}

	// Retrieving all data from cache
	cachedData := storage.ReadAll()

	// Verifying the number of entries matches
	if len(cachedData) != len(testData) {
		t.Errorf("Expected %d entries in cache, got %d", len(testData), len(cachedData))
	}

	// Verifying that all input data is present in the cache
	for _, expected := range testData {
		found := false
		for _, cached := range cachedData {
			if cached.Params.ObjectCost == expected.Params.ObjectCost {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected data with ObjectCost=%d not found in cache", expected.Params.ObjectCost)
		} else {
			t.Logf("Data with ObjectCost=%d found in cache", expected.Params.ObjectCost)
		}
	}
}

// TestHasData verifies that the HasData method correctly determines whether there is data in the cache.
func TestHasData(t *testing.T) {
	storage := cache.New()

	// Verifying empty cache
	if storage.HasData() {
		t.Error("Expected cache to be empty, but HasData() returned true")
	}

	// Adding data
	storage.Load(models.Result{Params: models.Params{ObjectCost: 100000}})

	// Verifying that cache is not empty
	if !storage.HasData() {
		t.Error("Expected cache to have data, but HasData() returned false")
	}
}

// TestConcurrentAccess verifies that the cache works correctly in a concurrent environment.
func TestConcurrentAccess(t *testing.T) {
	storage := cache.New()
	const numGoroutines = 100

	var wg sync.WaitGroup

	// Launching multiple goroutines that add data to the cache
	for range numGoroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			storage.Load(models.Result{Params: models.Params{ObjectCost: 100000}})
		}()
	}

	wg.Wait()

	// Verifying that all data has been added
	cachedData := storage.ReadAll()
	if len(cachedData) != numGoroutines {
		t.Errorf("Expected %d entries in cache, got %d", numGoroutines, len(cachedData))
	}

	// Verifying that IDCounter has incremented correctly
	if atomic.LoadInt32(&storage.IDCounter) != int32(numGoroutines) {
		t.Errorf("Expected IDCounter to be %d, got %d", numGoroutines, atomic.LoadInt32(&storage.IDCounter))
	}
}
