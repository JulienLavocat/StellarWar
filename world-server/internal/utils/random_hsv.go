package utils

import "math/rand"

func RandomHSV() [3]int {
	return [3]int{rand.Intn(360-1) + 1, rand.Intn(95-17) + 17, rand.Intn(100-80) + 80}
}
