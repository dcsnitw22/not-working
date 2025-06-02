package ipgenerator

import "fmt"

var RandomPduIP int = -1

func GetRandomPDUIP() string {

	if RandomPduIP < 255 {
		RandomPduIP = RandomPduIP + 1
	}

	randomPduIP := fmt.Sprintf("%d", RandomPduIP)

	ip := "172.31.20." + randomPduIP

	return ip

}
