package order

import (
	"crypto/rand"
	"fmt"
	"time"
)

const orderNumberAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func generateOrderNumber(location *time.Location) (string, error) {
	buffer := make([]byte, 6)
	if _, err := rand.Read(buffer); err != nil {
		return "", fmt.Errorf("generate order number: %w", err)
	}
	for index := range buffer {
		buffer[index] = orderNumberAlphabet[int(buffer[index])%len(orderNumberAlphabet)]
	}
	return fmt.Sprintf("JDC-%s-%s", time.Now().In(location).Format("20060102"), string(buffer)), nil
}
