package nas

import "strconv"

func DecodeHNPKI(hnpki byte) string {
	if int(hnpki) == 256 {
		return "Reserved"
	}
	val := "Home network PKI value" + strconv.Itoa(int(hnpki))
	return val
}
