package nas

import "errors"

func EncodeSecurityHeader(secHead string) (byte, error) {
	switch secHead {
	case "Plain 5GS NAS message":
		return PlainNAS, nil
	case "Integrity Protected":
		return IntegrityProtected, nil
	case "Integrity Protected and ciphered":
		return IntegrityCipher, nil
	case "Integrity protected with new 5G NAS security context":
		return IntegrityNew, nil
	case "Integrity protected and ciphered with new 5G NAS security context":
		return IntegrityCipherNew, nil
	default:
		return 0, errors.New("invalid Security Header Type")
	}
}
