package shared

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateAccountNumber() string {
	currentDate := time.Now().Format("20060102")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumber := fmt.Sprintf("%09d", r.Intn(1000000000))

	accountNumber := currentDate + randomNumber
	return accountNumber
}
