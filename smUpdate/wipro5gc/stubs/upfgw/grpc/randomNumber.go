package grpc

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

var RandomTunIP int = -1

func GetRandomID() string {

	// Seed the random number generator
	rand.Seed(uint64(time.Now().UnixNano()))
	// Generate a random number between 0 and 9999
	randomNumber := rand.Intn(10000) // Generates a number between 0 and 9999
	// Convert the number to a zero-padded 4-digit string
	randomString := fmt.Sprintf("%04d", randomNumber)

	return randomString

}

func GetRandomIP() string {

	if RandomTunIP < 255 {
		RandomTunIP = RandomTunIP + 1
	}

	randomIP := fmt.Sprintf("%d", RandomTunIP)

	ip := "172.24.20." + randomIP

	return ip

}
