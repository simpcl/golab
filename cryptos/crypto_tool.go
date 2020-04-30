package cryptos

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"errors"

	"github.com/simpcl/golabs/util"
)

type CipherType int32
type OpModeType int32
type PaddingType int32

const (
	Cipher_DES    CipherType  = 0
	Cipher_AES    CipherType  = 1
	OpMode_CBC    OpModeType  = 2
	OpMode_ECB    OpModeType  = 3
	Padding_No    PaddingType = 4
	Padding_PKCS5 PaddingType = 5
)

type CryptoTool struct {
	key     []byte
	cipher  CipherType
	opMode  OpModeType
	padding PaddingType
}

func NewCryptoTool(key []byte, cipher CipherType) *CryptoTool {
	return &CryptoTool{key: key, cipher: cipher, opMode: OpMode_CBC, padding: Padding_PKCS5}
}

func (ct *CryptoTool) SetOpModeType(om OpModeType) {
	ct.opMode = om
}

func (ct *CryptoTool) SetPaddingType(pt PaddingType) {
	ct.padding = pt
}

func (ct *CryptoTool) Encrypt(data []byte) ([]byte, error) {
	var block cipher.Block
	var err error

	if ct.cipher == Cipher_DES {
		block, err = des.NewCipher(ct.key)
	} else if ct.cipher == Cipher_AES {
		block, err = aes.NewCipher(ct.key)
	} else {
		return nil, errors.New("Not supported cipher type")
	}
	if err != nil {
		return nil, err
	}

	src := []byte(data)
	if ct.padding == Padding_PKCS5 {
		src = util.PKCS5Padding(src, block.BlockSize())
	}
	dst := make([]byte, len(src))

	var blockMode cipher.BlockMode
	if ct.opMode == OpMode_ECB {
		blockMode = util.NewECBEncrypter(block)
	} else if ct.opMode == OpMode_CBC {
		blockMode = cipher.NewCBCEncrypter(block, bytes.Repeat([]byte("0"), block.BlockSize()))
	} else {
		return nil, errors.New("Not supported operation mode")
	}
	blockMode.CryptBlocks(dst, src)

	return dst, nil
}

func (ct *CryptoTool) Decrypt(encrytpedData []byte) ([]byte, error) {
	var block cipher.Block
	var err error

	if ct.cipher == Cipher_DES {
		block, err = des.NewCipher(ct.key)
	} else if ct.cipher == Cipher_AES {
		block, err = aes.NewCipher(ct.key)
	} else {
		return nil, errors.New("Not supported cipher type")
	}
	if err != nil {
		return nil, err
	}

	dst := make([]byte, len(encrytpedData))

	var blockMode cipher.BlockMode
	if ct.opMode == OpMode_ECB {
		blockMode = util.NewECBDecrypter(block)
	} else if ct.opMode == OpMode_CBC {
		blockMode = cipher.NewCBCDecrypter(block, bytes.Repeat([]byte("0"), block.BlockSize()))
	} else {
		return nil, errors.New("Not supported operation mode")
	}

	blockMode.CryptBlocks(dst, encrytpedData)
	if ct.padding == Padding_PKCS5 {
		dst = util.UnPKCS5Padding(dst)
	}

	return dst, nil
}
