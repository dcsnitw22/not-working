package nas

import (
	"errors"
)

//Logic to Encode PDU Session Type

func EncodePduSessionType(pdustype string) (byte, error) {
	switch pdustype {
	case "IPV4":
		return IpV4, nil
	case "IPV6":
		return IpV6, nil
	case "IPV4V6":
		return IPV4V6, nil
	case "UNSTRUCTURED":
		return Unstructured, nil
	case "ETHERNET":
		return Ethernet, nil
	default:
		return 0, errors.New("invalid PDU Session Type")
	}
}
