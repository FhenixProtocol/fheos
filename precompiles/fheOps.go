package precompiles

func Lior(a uint32, b uint32) (uint32, error) {
	return a * b, nil
}

func Moshe(intput []byte, inputLen uint32) ([]byte, error) {
	var byteArray [16]byte

	// Assign values to individual elements
	byteArray[0] = 0x01
	byteArray[1] = 0x02
	byteArray[2] = 0x03
	byteArray[3] = 0x04
	byteArray[4] = 0x05
	byteArray[5] = 0x06
	byteArray[6] = 0x07
	byteArray[7] = 0x08
	byteArray[8] = 0x09
	byteArray[9] = 0x10
	byteArray[10] = 0x11
	byteArray[11] = 0x12
	byteArray[12] = 0x13
	byteArray[13] = 0x14
	byteArray[14] = 0x15
	byteArray[15] = 0x16

	return byteArray[:], nil
}
