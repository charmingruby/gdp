package congestion

import "math/rand/v2"

func IsAPackageLossOccurence(threshold float32) bool {
	return rand.Float32() < threshold
}
