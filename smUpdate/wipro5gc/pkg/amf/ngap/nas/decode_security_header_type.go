package nas

import "errors"

func DecodeSecurityHeader(secHead byte) (string, error) {
	switch secHead {
	case PlainNAS:
		return "Plain 5GS NAS message", nil
	case IntegrityProtected:
		return "Integrity Protected", nil
	case IntegrityCipher:
		return "Integrity Protected and ciphered", nil
	case IntegrityNew:
		return "Integrity protected with new 5G NAS security context", nil
	case IntegrityCipherNew:
		return "Integrity protected and ciphered with new 5G NAS security context", nil
	default:
		return "", errors.New("invalid Security Header Type")
	}
}
