package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/aliceKatkout/test_fizzbuzz/internal/fizzbuzz"
	"github.com/aliceKatkout/test_fizzbuzz/internal/stats"
)

func parseFizzBuzzRequest(q url.Values) (fizzbuzz.Request, error) {
	int1, err := parseIntParam(q, "int1")
	if err != nil {
		return fizzbuzz.Request{}, err
	}
	int2, err := parseIntParam(q, "int2")
	if err != nil {
		return fizzbuzz.Request{}, err
	}
	limit, err := parseIntParam(q, "limit")
	if err != nil {
		return fizzbuzz.Request{}, err
	}

	req := fizzbuzz.Request{
		Int1:  int1,
		Int2:  int2,
		Limit: limit,
		Str1:  q.Get("str1"),
		Str2:  q.Get("str2"),
	}
	return req, nil
}

func parseIntParam(q url.Values, name string) (int, error) {
	raw := q.Get(name)
	if raw == "" {
		return 0, fmt.Errorf("%s is required", name)
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer", name)
	}
	return v, nil
}

// FizzBuzzHandler returns the handler for GET /fizzbuzz. Every valid request
// is recorded in store for the statistics endpoint.
func FizzBuzzHandler(store *stats.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := parseFizzBuzzRequest(r.URL.Query())
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := req.Validate(); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		store.Record(stats.Key{
			Int1:  req.Int1,
			Int2:  req.Int2,
			Limit: req.Limit,
			Str1:  req.Str1,
			Str2:  req.Str2,
		})

		writeJSON(w, http.StatusOK, fizzBuzzResponse{Result: fizzbuzz.GenerateFizzbuzz(req)})
	}
}

type fizzBuzzResponse struct {
	Result []string `json:"result"`
}

// StatisticsHandler returns the handler for GET /statistics: the most
// frequently requested FizzBuzz parameters and their hit count.
func StatisticsHandler(store *stats.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key, hits, ok := store.Top()
		if !ok {
			writeJSON(w, http.StatusOK, statisticsResponse{Request: nil, Hits: 0})
			return
		}

		writeJSON(w, http.StatusOK, statisticsResponse{
			Request: &statisticsRequest{
				Int1:  key.Int1,
				Int2:  key.Int2,
				Limit: key.Limit,
				Str1:  key.Str1,
				Str2:  key.Str2,
			},
			Hits: hits,
		})
	}
}

type statisticsRequest struct {
	Int1  int    `json:"int1"`
	Int2  int    `json:"int2"`
	Limit int    `json:"limit"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
}

type statisticsResponse struct {
	Request *statisticsRequest `json:"request"`
	Hits    int                `json:"hits"`
}
