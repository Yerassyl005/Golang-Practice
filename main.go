package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

func IsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		return true // network error
	}

	if resp == nil {
		return false
	}

	switch resp.StatusCode {
	case 429, 500, 502, 503, 504:
		return true
	case 401, 404:
		return false
	default:
		return false
	}
}

func CalculateBackoff(attempt int) time.Duration {
	baseDelay := 500 * time.Millisecond
	maxDelay := 5 * time.Second

	backoff := baseDelay * time.Duration(math.Pow(2, float64(attempt)))

	if backoff > maxDelay {
		backoff = maxDelay
	}

	jitter := time.Duration(rand.Int63n(int64(backoff)))
	return jitter
}

func ExecutePayment(ctx context.Context, url string) error {
	client := &http.Client{}

	maxRetries := 5

	for attempt := 0; attempt < maxRetries; attempt++ {

		// Проверка context
		if ctx.Err() != nil {
			return ctx.Err()
		}

		req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

		resp, err := client.Do(req)

		if err == nil && resp.StatusCode == 200 {
			fmt.Println("Attempt", attempt+1, ": SUCCESS")
			return nil
		}

		if !IsRetryable(resp, err) {
			return fmt.Errorf("non-retryable error")
		}

		if attempt == maxRetries-1 {
			break
		}

		delay := CalculateBackoff(attempt)

		fmt.Printf("Attempt %d failed, waiting %v...\n", attempt+1, delay)

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("failed after retries")
}


type CachedResponse struct {
	StatusCode int
	Body       []byte
	Completed  bool
}

type MemoryStore struct {
	mu   sync.Mutex
	data map[string]*CachedResponse
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]*CachedResponse),
	}
}

func (m *MemoryStore) Get(key string) (*CachedResponse, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	resp, ok := m.data[key]
	return resp, ok
}

func (m *MemoryStore) StartProcessing(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[key]; exists {
		return false
	}

	m.data[key] = &CachedResponse{Completed: false}
	return true
}

func (m *MemoryStore) Finish(key string, status int, body []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	resp := m.data[key]
	resp.StatusCode = status
	resp.Body = body
	resp.Completed = true
}

func IdempotencyMiddleware(store *MemoryStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		key := r.Header.Get("Idempotency-Key")

		if key == "" {
			http.Error(w, "Idempotency-Key required", http.StatusBadRequest)
			return
		}

		if cached, exists := store.Get(key); exists {
			if cached.Completed {
				w.WriteHeader(cached.StatusCode)
				w.Write(cached.Body)
			} else {
				http.Error(w, "Processing", http.StatusConflict)
			}
			return
		}

		if !store.StartProcessing(key) {
			http.Error(w, "Conflict", http.StatusConflict)
			return
		}

		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		store.Finish(key, rec.Code, rec.Body.Bytes())

		w.WriteHeader(rec.Code)
		w.Write(rec.Body.Bytes())
	})
}


func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("===== TASK 1 =====")

	counter := 0

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter++

		if counter <= 3 {
			w.WriteHeader(503)
			return
		}

		w.WriteHeader(200)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer testServer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := ExecutePayment(ctx, testServer.URL)
	if err != nil {
		fmt.Println("Final error:", err)
	}


	fmt.Println("\n===== TASK 2 =====")

	store := NewMemoryStore()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Processing started")

		time.Sleep(2 * time.Second)

		resp := map[string]interface{}{
			"status":         "paid",
			"amount":         1000,
			"transaction_id": "uuid-123",
		}

		json.NewEncoder(w).Encode(resp)
	})

	server := httptest.NewServer(IdempotencyMiddleware(store, handler))
	defer server.Close()

	wg := sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			req, _ := http.NewRequest("POST", server.URL, nil)
			req.Header.Set("Idempotency-Key", "same-key")

			client := &http.Client{}
			resp, err := client.Do(req)

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			fmt.Println("Response status:", resp.StatusCode)
		}(i)
	}

	wg.Wait()
}
