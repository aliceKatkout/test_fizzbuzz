package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aliceKatkout/test_fizzbuzz/internal/stats"
)

func TestFizzBuzzHandler_Success(t *testing.T) {
	store := stats.NewStore()
	handler := FizzBuzzHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/fizzbuzz?int1=3&int2=5&limit=15&str1=fizz&str2=buzz", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var got fizzBuzzResponse
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	want := strings.Split("1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz", ",")
	if len(got.Result) != len(want) {
		t.Fatalf("got %d items, want %d", len(got.Result), len(want))
	}
	for i := range want {
		if got.Result[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got.Result[i], want[i])
		}
	}
}

func TestFizzBuzzHandler_RecordsStats(t *testing.T) {
	store := stats.NewStore()
	handler := FizzBuzzHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/fizzbuzz?int1=3&int2=5&limit=15&str1=fizz&str2=buzz", nil)
	handler(httptest.NewRecorder(), req)

	key, hits, ok := store.Top()
	if !ok {
		t.Fatal("expected a recorded request")
	}
	want := stats.Key{Int1: 3, Int2: 5, Limit: 15, Str1: "fizz", Str2: "buzz"}
	if key != want || hits != 1 {
		t.Errorf("got key=%+v hits=%d, want key=%+v hits=1", key, hits, want)
	}
}

func TestFizzBuzzHandler_InvalidParams(t *testing.T) {
	tests := []struct {
		name  string
		query string
	}{
		{"missing int1", "int2=5&limit=15&str1=fizz&str2=buzz"},
		{"non-integer int1", "int1=abc&int2=5&limit=15&str1=fizz&str2=buzz"},
		{"zero int1", "int1=0&int2=5&limit=15&str1=fizz&str2=buzz"},
		{"negative limit", "int1=3&int2=5&limit=-1&str1=fizz&str2=buzz"},
		{"limit too large", "int1=3&int2=5&limit=99999999&str1=fizz&str2=buzz"},
		{"empty str1", "int1=3&int2=5&limit=15&str1=&str2=buzz"},
	}

	store := stats.NewStore()
	handler := FizzBuzzHandler(store)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/fizzbuzz?"+tt.query, nil)
			rec := httptest.NewRecorder()
			handler(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestStatisticsHandler_Empty(t *testing.T) {
	store := stats.NewStore()
	handler := StatisticsHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/statistics", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var got statisticsResponse
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.Request != nil || got.Hits != 0 {
		t.Errorf("got %+v, want empty response", got)
	}
}

func TestStatisticsHandler_ReturnsMostFrequent(t *testing.T) {
	store := stats.NewStore()
	fizzBuzzHandler := FizzBuzzHandler(store)

	makeRequest := func(query string) {
		req := httptest.NewRequest(http.MethodGet, "/fizzbuzz?"+query, nil)
		fizzBuzzHandler(httptest.NewRecorder(), req)
	}

	makeRequest("int1=3&int2=5&limit=15&str1=fizz&str2=buzz")
	makeRequest("int1=3&int2=5&limit=15&str1=fizz&str2=buzz")
	makeRequest("int1=2&int2=7&limit=50&str1=foo&str2=bar")

	rec := httptest.NewRecorder()
	StatisticsHandler(store)(rec, httptest.NewRequest(http.MethodGet, "/statistics", nil))

	var got statisticsResponse
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got.Hits != 2 {
		t.Errorf("got hits=%d, want 2", got.Hits)
	}
	if got.Request == nil || got.Request.Int1 != 3 || got.Request.Int2 != 5 || got.Request.Limit != 15 ||
		got.Request.Str1 != "fizz" || got.Request.Str2 != "buzz" {
		t.Errorf("got request=%+v, want the {3,5,15,fizz,buzz} request", got.Request)
	}
}

func TestRouter_HealthCheck(t *testing.T) {
	router := NewRouter(stats.NewStore())
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}
}
