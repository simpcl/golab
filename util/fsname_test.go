package util

import (
	"encoding/binary"
	"math/rand"
	"testing"
	"time"
)

func TestXorKeyMask(t *testing.T) {
	fsnc := NewFsnameCoder("CloudStorage")

	var input = []byte("abcdefghijklmnopqrstuvwxyz")
	var output = make([]byte, len(input))
	fsnc.xorKeyMask(input, output)
	var output2 = make([]byte, len(input))
	fsnc.xorKeyMask(output, output2)
	for i := 0; i < len(input); i++ {
		if input[i] != output2[i] {
			t.Errorf("not equal, i: %d", i)
			return
		}
	}
}

func TestEncodeDecode(t *testing.T) {
	fsnc := NewFsnameCoder("CloudStorage")

	var input = make([]byte, 12)
	var output = make([]byte, 16)
	var output2 = make([]byte, 12)

	for i := 0; i < 10; i++ {
		rand.Seed(time.Now().UnixNano())
		binary.BigEndian.PutUint32(input, uint32(rand.Int31()))
		binary.BigEndian.PutUint32(input[4:], uint32(rand.Int31()))
		binary.BigEndian.PutUint32(input[8:], uint32(rand.Int31()))
		t.Logf("encode input: %x\n", input)
		if err := fsnc.encode(input, output); err != nil {
			t.Errorf("encode error: %v", err)
			return
		}
		t.Logf("encode output: %s\n", string(output))

		if err := fsnc.decode(output, output2); err != nil {
			t.Errorf("decode error: %v", err)
			return
		}
		t.Logf("decode output: %x\n", output2)

		for i := 0; i < len(input); i++ {
			if input[i] != output2[i] {
				t.Errorf("not equal, i: %d", i)
				return
			}
		}
	}
}

func TestFname(t *testing.T) {
	fsnc := NewFsnameCoder("CloudStorage")

	for i := 0; i < 10; i++ {
		appId := uint32(1)
		fileId := uint64(i + 1)

		fname, err := fsnc.EncodeName(appId, fileId)
		if err != nil {
			t.Errorf("GenerateFname Error: %v", err)
			return
		}
		t.Logf("appId: %d, fileId: %d, fname: %s\n", appId, fileId, fname)

		appId2, fileId2, err2 := fsnc.DecodeName(fname)
		if err2 != nil {
			t.Errorf("ParseFname Error: %v", err2)
			return
		}

		if appId != appId2 || fileId != fileId2 {
			t.Errorf("appId or fileId are not equal: %d <-> %d, %d <-> %d", appId, appId2, fileId, fileId2)
			return
		}
	}
}
