package pgconn

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/tjfoc/gmsm/sm3"
	"golang.org/x/crypto/pbkdf2"
)

func RFC5802Algorithm(password string, random64code string, token string, serverSignature string, serverIteration int, method string) (string, error) {

	k := generateKFromPBKDF2(password, random64code, serverIteration)

	serverKey := computeHMAC(k, []byte("Sever Key"))
	clientKey := computeHMAC(k, []byte("Client Key"))
	var storedKey []byte

	if strings.EqualFold(method, "sha256") {
		storedKey = getSha256(clientKey)
	} else if strings.EqualFold(method, "sm3") {
		storedKey = sm3.Sm3Sum(clientKey)
	}

	tokenByte, err := hex.DecodeString(token)
	if err != nil {
		return "", err
	}

	clientSignature := computeHMAC(serverKey, tokenByte)
	if serverSignature != "" && serverSignature != hex.EncodeToString(clientSignature) {
		return "", fmt.Errorf("serverSignature(%s) != clientSignature(%s)", serverSignature, hex.EncodeToString(clientSignature))
	}
	hmacResult := computeHMAC(storedKey, tokenByte)
	h := XorBetweenPassword(hmacResult, clientKey, len(clientKey))
	
	return hex.EncodeToString(h), nil

}

func generateKFromPBKDF2(password string, random64code string, serverIteration int) []byte {
	random32code, err := hex.DecodeString(random64code)
	if err != nil {
		fmt.Println(err.Error())
	}
	pwdEn := pbkdf2.Key([]byte(password), random32code, serverIteration, 32, sha1.New)
	return pwdEn
}
func getSha256(message []byte) []byte {
	hash := sha256.New()
	hash.Write(message)

	return hash.Sum(nil)
}
func XorBetweenPassword(password1 []byte, password2 []byte, length int) []byte {
	array := make([]byte, length)
	for i := 0; i < length; i++ {
		array[i] = (password1[i] ^ password2[i])
	}
	return array
}
