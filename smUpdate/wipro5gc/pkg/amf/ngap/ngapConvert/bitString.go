package ngapConvert

import (
	"encoding/hex"
	"fmt"

	"w5gc.io/wipro5gcore/pkg/amf/ngap/asn1gen/asn1rt"
)

// func BitStringToHex(bitString *aper.BitString) (hexString string) {
func BitStringToHex(bitString *asn1rt.BitString) (hexString string) {
	hexString = hex.EncodeToString(bitString.Bytes)
	hexLen := (bitString.BitLength + 3) / 4
	hexString = hexString[:hexLen]
	return
}

func HexToBitString(hexString string, bitLength int) (bitString asn1rt.BitString) {
	hexLen := len(hexString)
	if hexLen != (bitLength+3)/4 {
		fmt.Println("hexLen[", hexLen, "] doesn't match bitLength[", bitLength, "]")
		return
	}
	if hexLen%2 == 1 {
		hexString += "0"
	}
	if byteTmp, err := hex.DecodeString(hexString); err != nil {
		fmt.Printf("Decode byteString failed: %+v", err)
	} else {
		bitString.Bytes = byteTmp
	}
	bitString.BitLength = int(bitLength)
	mask := byte(0xff)
	mask = mask << uint(8-bitLength%8)
	if mask != 0 {
		bitString.Bytes[len(bitString.Bytes)-1] &= mask
	}
	return
}
