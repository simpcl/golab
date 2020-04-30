package cryptos

import "testing"

func Test_DES_CBC_NoPadding_Decrypt(t *testing.T) {
	desTool := NewCryptoTool([]byte("12345678"), Cipher_DES)
	desTool.SetOpModeType(OpMode_CBC)
	desTool.SetPaddingType(Padding_No)

	src := "DES_CBC_NoPadding_String"
	encryptedData, err := desTool.Encrypt([]byte(src))
	if err != nil {
		t.Errorf("des encrypt: %s\n", err.Error())
		return
	}
	dst, err := desTool.Decrypt(encryptedData)
	if err != nil {
		t.Errorf("des decrypt: %s\n", err.Error())
		return
	}
	if src != string(dst) {
		t.Errorf("not equal: %s %s\n", src, string(dst))
	}
}

func Test_DES_ECB_PKSC5Padding_Decrypt(t *testing.T) {
	desTool := NewCryptoTool([]byte("12345678"), Cipher_DES)
	desTool.SetOpModeType(OpMode_ECB)
	desTool.SetPaddingType(Padding_PKCS5)

	src := "DES_ECB_PKSC5Padding_String"
	encryptedData, err := desTool.Encrypt([]byte(src))
	if err != nil {
		t.Errorf("des encrypt: %s\n", err.Error())
		return
	}
	dst, err := desTool.Decrypt(encryptedData)
	if err != nil {
		t.Errorf("des decrypt: %s\n", err.Error())
		return
	}
	if src != string(dst) {
		t.Errorf("not equal: %s %s\n", src, string(dst))
	}
}

func Test_AES_CBC_PKSC5Padding_Decrypt(t *testing.T) {
	aesTool := NewCryptoTool([]byte("0123456789abcdef"), Cipher_AES)
	aesTool.SetOpModeType(OpMode_CBC)
	aesTool.SetPaddingType(Padding_PKCS5)

	src := "AES_CBC_NoPadding_String"
	encryptedData, err := aesTool.Encrypt([]byte(src))
	if err != nil {
		t.Errorf("aes encrypt: %s\n", err.Error())
		return
	}
	dst, err := aesTool.Decrypt(encryptedData)
	if err != nil {
		t.Errorf("aes decrypt: %s\n", err.Error())
		return
	}
	if src != string(dst) {
		t.Errorf("not equal: %s %s\n", src, string(dst))
	}
}

func Test_AES_ECB_PKSC5Padding_Decrypt(t *testing.T) {
	aesTool := NewCryptoTool([]byte("0123456789abcdef"), Cipher_AES)
	aesTool.SetOpModeType(OpMode_ECB)
	aesTool.SetPaddingType(Padding_PKCS5)

	src := "AES_CBC_PKSC5Padding_String"
	encryptedData, err := aesTool.Encrypt([]byte(src))
	if err != nil {
		t.Errorf("aes encrypt: %s\n", err.Error())
		return
	}
	dst, err := aesTool.Decrypt(encryptedData)
	if err != nil {
		t.Errorf("aes decrypt: %s\n", err.Error())
		return
	}
	if src != string(dst) {
		t.Errorf("not equal: %s %s\n", src, string(dst))
	}
}
