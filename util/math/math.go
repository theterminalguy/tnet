package mathutil

import "math/rand"

func RandomFloat32(n []float32) float32 {
	if len(n) == 0 {
		n = []float32{1.0, 2.5, 3, 4.0, 5.2, 8.5, 90, 10}
		return n[rand.Intn(len(n))]
	}
	return float32(rand.Float32())
}
