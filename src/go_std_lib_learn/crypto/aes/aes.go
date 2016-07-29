package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// NewCBCDecrypter 使用AES(CBC)模式解密
func NewCBCDecrypter(inKey []byte, inCiphertext []byte, inIV []byte) (outPlanText []byte) {
	key := inKey
	ciphertext := inCiphertext
	iv := inIV
	/*
		fmt.Printf("====\n")
		fmt.Printf("key=%s\n", key)
		fmt.Printf("ciphertext=%x\n", ciphertext)
		fmt.Printf("iv=%x\n", iv)
		fmt.Printf("====\n")
	*/
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	return ciphertext
}

// NewCBCEncrypter 使用AES(CBC)模式加密
func NewCBCEncrypter(inKey []byte, inPlantext []byte) (outCiphertext []byte, usedIV []byte) {
	key := inKey
	plaintext := inPlantext

	// CBC mode works on blocks so plaintexts may need to be padded to the
	// next whole block. For an example of such padding, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. Here we'll
	// assume that the plaintext is already of the correct length.
	if len(plaintext)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(plaintext, plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	//fmt.Printf("ciphertext = %x\n", ciphertext)
	return plaintext, iv
}

func translationBinaryToChar(inBinary []byte) (outChar []byte) {
	if inBinary == nil || len(inBinary) == 0 {
		fmt.Println("input binary is NULL")
		return nil
	}

	j := 0
	oChar := make([]byte, 32)
	for _, value := range inBinary {
		oChar[j] = value&0x0f + 'A'
		j++
		oChar[j] = (value>>4)&0x0f + 'A'
		j++
	}
	return oChar
}

func main() {
	key := []byte("CBSV5R5SrvPubKey")
	planText := []byte("1234567890123456")

	oCiphertext, uIV := NewCBCEncrypter(key, planText)
	fmt.Printf("ciphertext = %x\nUsed IV = %x\n", oCiphertext, uIV)
	outstr := translationBinaryToChar(oCiphertext)
	fmt.Printf("B2C = %x\n", outstr)
	/*
		fmt.Printf("key=%s\n", key)
		fmt.Printf("ciphertext=%x\n", oCiphertext)
		fmt.Printf("iv=%x\n", uIV)
	*/
	oChar := NewCBCDecrypter(key, oCiphertext, uIV)
	fmt.Printf("planText = %s\n", oChar)
}
