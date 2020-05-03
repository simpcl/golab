package util

import (
	"encoding/binary"
	"errors"
	"math/rand"
	"time"
)

type FsnameCoder struct {
	keyMask       []byte
	keyMaskLength int
	encodeTable   []byte
	decodeTable   []uint8
}

func NewFsnameCoder(key string) *FsnameCoder {
	coder := &FsnameCoder{}
	coder.keyMask = []byte(key)
	coder.keyMaskLength = len(coder.keyMask)
	coder.encodeTable = []byte("0JoU8EaN3xf19hIS2d.6pZRFBYurMDGw7K5m4CyXsbQjg_vTOAkcHVtzqWilnLPe")
	coder.decodeTable = []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 18, 0, 0, 11, 16, 8, 36, 34, 19, 32, 4, 12, 0, 0, 0, 0, 0, 0, 0, 49, 24, 37, 29, 5, 23, 30, 52, 14, 1, 33, 61, 28, 7, 48, 62, 42, 22, 15, 47, 3, 53, 57, 39, 25, 21, 0, 0, 0, 0, 45, 0, 6, 41, 51, 17, 63, 10, 44, 13, 58, 43, 50, 59, 35, 60, 2, 20, 56, 27, 40, 54, 26, 46, 31, 9, 38, 55, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	return coder
}

func (fsnc *FsnameCoder) xorKeyMask(input []byte, output []byte) {
	for i, _ := range input {
		output[i] = input[i] ^ fsnc.keyMask[i%fsnc.keyMaskLength]
	}
}

func (fsnc *FsnameCoder) encode(input []byte, output []byte) error {
	if input == nil || len(input) <= 0 {
		return errors.New("invalid input")
	}
	inputLength := len(input)
	if inputLength%3 != 0 {
		return errors.New("invalid input length")
	}
	if output == nil || len(output) <= 0 {
		return errors.New("invalid output")
	}

	buffer := make([]byte, inputLength)
	fsnc.xorKeyMask(input, buffer)

	var k int32 = 0
	for i := 0; i < inputLength; i += 3 {
		value := ((uint32(buffer[i]) << 16) & 0xff0000) + ((uint32(buffer[i+1]) << 8) & 0xff00) + (uint32(buffer[i+2]) & 0xff)
		output[k] = fsnc.encodeTable[value>>18]
		k += 1
		output[k] = fsnc.encodeTable[(value>>12)&0x3f]
		k += 1
		output[k] = fsnc.encodeTable[(value>>6)&0x3f]
		k += 1
		output[k] = fsnc.encodeTable[value&0x3f]
		k += 1
	}
	return nil
}

func (fsnc *FsnameCoder) decode(input []byte, output []byte) error {
	if input == nil || len(input) <= 0 {
		return errors.New("invalid input")
	}
	inputLength := len(input)
	if inputLength%4 != 0 {
		return errors.New("invalid input length")
	}
	if output == nil || len(output) <= 0 {
		return errors.New("invalid output")
	}

	bufferLength := inputLength / 4 * 3
	buffer := make([]byte, bufferLength)
	var k int32 = 0
	for i := 0; i < inputLength; i += 4 {
		value := uint32(fsnc.decodeTable[int(input[i])])<<18 + uint32(fsnc.decodeTable[int(input[i+1])])<<12 + uint32(fsnc.decodeTable[int(input[i+2])])<<6 + uint32(fsnc.decodeTable[int(input[i+3])])
		buffer[k] = byte((value >> 16) & 0xff)
		k++
		buffer[k] = byte((value >> 8) & 0xff)
		k++
		buffer[k] = byte((value) & 0xff)
		k++
	}
	fsnc.xorKeyMask(buffer, output)
	return nil
}

func (fsnc *FsnameCoder) EncodeName(appId uint32, fileId uint64) (string, error) {
	var input = make([]byte, 18)
	var output = make([]byte, 24)

	rand.Seed(time.Now().UnixNano())
	rn := rand.Int63()
	binary.BigEndian.PutUint64(input, uint64(rn))
	binary.BigEndian.PutUint32(input[6:], uint32(appId))
	binary.BigEndian.PutUint64(input[10:], uint64(fileId))

	if err := fsnc.encode(input, output); err != nil {
		return "", errors.New("fsname encode error: " + err.Error())
	}
	return string(output), nil
}

func (fsnc *FsnameCoder) DecodeName(fname string) (uint32, uint64, error) {
	input := []byte(fname)
	output := make([]byte, len(input))
	if err := fsnc.decode(input, output); err != nil {
		return 0, 0, errors.New("fsname decode error: " + err.Error())
	}
	appId := binary.BigEndian.Uint32(output[6:])
	fileId := binary.BigEndian.Uint64(output[10:])
	return appId, fileId, nil
}
