package udp

import "math/rand/v2"

type Threshold struct {
	PackageLoss float32
}

func isAPackageLossOccurence(threshold float32) bool {
	return rand.Float32() < threshold
}
