package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

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
}
