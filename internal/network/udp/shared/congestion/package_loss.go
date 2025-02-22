package congestion

import "math/rand/v2"

func (c *Congestion) IsAPackageLossOccurence(threshold float32) bool {
	return rand.Float32() < threshold
}
