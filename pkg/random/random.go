package random

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"
)

// Token generates a secure random hex token
func Token(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// OrderNumber generates unique order number: ZAEM-YYYYMMDD-XXXXX
func OrderNumber() string {
	date := time.Now().Format("20060102")
	randomPart := randomString(5, "ABCDEFGHJKLMNPQRSTUVWXYZ23456789")
	return fmt.Sprintf("ZAEM-%s-%s", date, randomPart)
}

// NumericCode generates a numeric OTP
func NumericCode(length int) string {
	code := ""
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		code += fmt.Sprintf("%d", n.Int64())
	}
	return code
}

func randomString(length int, charset string) string {
	result := make([]byte, length)
	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[n.Int64()]
	}
	return string(result)
}