package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

func NewCBCDecrypter() {
	key := []byte("example key 1234")
	ciphertext, _ := hex.DecodeString("f363f3ccdcb12bb883abf484ba77d9cd7d32b5baecb3d4b1b3e0e4beffdb3ded")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	// If the original plaintext lengths are not a multiple of the block
	// size, padding would have to be added when encrypting, which would be
	// removed at this point. For an example, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. However, it's
	// critical to note that ciphertexts must be authenticated (i.e. by
	// using crypto/hmac) before being decrypted in order to avoid creating
	// a padding oracle.

	fmt.Printf("ciphertext = %s\n", ciphertext)
	// Output: exampleplaintext
}

func NewCBCEncrypter(inKey []byte, inPlantext string) (outCiphertext []byte, usedIV []byte) {
	key := inKey
	plaintext := []byte(inPlantext)

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
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	//fmt.Printf("ciphertext = %x\n", ciphertext)
	return ciphertext, iv
}

func translationBinaryToChar(inBinary []byte) (outChar []byte) {
	if inBinary == nil || len(inBinary) == 0 {
		fmt.Println("input binary is NULL")
		return nil
	}

	j := 0
	oChar := make([]byte, 68)
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

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte("108000000000000001|2000|2114380800"))
	//expectedMAC := mac.Sum([]byte("1"))
	expectedMAC := mac.Sum(nil)
	fmt.Printf("expectedMAC = %v\n", expectedMAC)
	fmt.Printf("expectedMAC = %v\n", string(expectedMAC))
	fmt.Printf("expectedMAC = %v\n", hex.EncodeToString(expectedMAC))
	outstr := translationBinaryToChar(expectedMAC)
	fmt.Printf("expectedMAC = %v\n", string(outstr))

	oCiphertext, uIV := NewCBCEncrypter(key, "1234567890123456")
	fmt.Printf("expectedAES = %x\nUsed IV = %x\n", oCiphertext, uIV)
	NewCBCDecrypter()
}
