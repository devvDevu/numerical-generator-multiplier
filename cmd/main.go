package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Config struct {
	RTP float64
}

var cfg Config

func main() {
	rtp := flag.Float64("rtp", 0.95, "Target RTP value (0 < rtp <= 1.0)")
	flag.Parse()

	if *rtp <= 0 || *rtp > 1.0 {
		log.Fatal("RTP must be in (0, 1.0]")
	}

	cfg = Config{RTP: *rtp}
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/get", handleGet)
	log.Println("Server started on :64333")
	log.Fatal(http.ListenAndServe(":64333", nil))
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	x, err := parseRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	multiplier := generateMultiplier(x)
	sendResponse(w, multiplier)
}

func parseRequest(r *http.Request) (float64, error) {
	xStr := r.URL.Query().Get("x")
	if xStr == "" {
		return 0, &ApiError{Message: "Parameter 'x' is required"}
	}

	x, err := strconv.ParseFloat(xStr, 64)
	if err != nil || x < 1.0 || x > 10000.0 {
		return 0, &ApiError{Message: "Invalid x value. Must be 1.0 <= x <= 10000.0"}
	}

	return x, nil
}

// generate a random multiplier by rtp
// Example:
// x = 1000
// rtp = 0.95
// if rand.Float64() < 0.95, then multiplier = x + rand.Float64()*(10000.0-x): 1000 + 0.05*(10000-1000) = 1450
// else multiplier = 1.0 + rand.Float64()*(x-1.0): 1.0 + 0.05*(1000-1.0) = 500,95

func generateMultiplier(x float64) float64 {
	if rand.Float64() < cfg.RTP {
		return x + rand.Float64()*(10000.0-x)
	}
	return 1.0 + rand.Float64()*(x-1.0)
}

func sendResponse(w http.ResponseWriter, multiplier float64) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{
		"result": multiplier,
	})
}

type ApiError struct {
	Message string `json:"error"`
}

func (e *ApiError) Error() string {
	return e.Message
}
