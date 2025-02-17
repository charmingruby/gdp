package packet

const (
	ackIDPacketSize  = 4
	sequentialIDSize = 4
	dataSize         = 1024
)

func DataSize() int {
	return dataSize
}
