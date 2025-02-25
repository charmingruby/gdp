package congestion

type Congestion struct {
	PackageLoss float32
	Cwnd        uint32
	Sshthresh   uint32
}
