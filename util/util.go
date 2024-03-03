package util

import (
	"math/rand"
	"time"
)

func CreateNRandomVectorsDimM(n, m int) [][]float32 {
	rand.Seed(time.Now().UnixNano())

	// Create and populate the vectors
	vectors := make([][]float32, n)
	for i := 0; i < n; i++ {
		vector := make([]float32, m)
		for j := 0; j < m; j++ {
			// Generate random values between 0 and 1
			vector[j] = rand.Float32()
		}
		vectors[i] = vector
	}

	// Print the generated vectors
	//for i, vector := range vectors {
	//	fmt.Printf("Vector %d: %v\n", i+1, vector)
	//}

	return vectors
}
