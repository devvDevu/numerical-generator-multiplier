package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	sequence := generateSequence(10000000)
	// sum of sequence
	sum0 := calculateSum(sequence)
	// sum of sequence, where multiplier > x
	sum1 := 0.0

	massMult := make([]float64, 0)

	for _, x := range sequence {
		multiplier := getMultiplier(x)
		massMult = append(massMult, multiplier)
		if multiplier > x {
			sum1 += x
		}
	}

	fmt.Printf("RTP: %.1f\n", sum1/sum0)
}

func calculateSum(seq []float64) float64 {
	sum := 0.0
	for _, x := range seq {
		sum += x
	}
	return sum
}

func generateSequence(n int) []float64 {
	seq := make([]float64, n)
	for i := range seq {
		seq[i] = 1.0 + rand.Float64()*9999.0
	}
	return seq
}

func getMultiplier(x float64) float64 {
	resp, _ := http.Get(fmt.Sprintf("http://localhost:64333/get?x=%.2f", x))
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct{ Result float64 }
	json.Unmarshal(body, &result)
	return result.Result
}
