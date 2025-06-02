package ngapConvert

import (
	"encoding/hex"
	"fmt"

	"w5gc.io/wipro5gcore/asn1gen/asn1rt"
)

func AmfIdToNgap(amfId string) (regionId, setId, ptrId asn1rt.BitString) {
	regionId = HexToBitString(amfId[:2], 8)
	setId = HexToBitString(amfId[2:5], 10)
	tmpByte, err := hex.DecodeString(amfId[4:])
	if err != nil {
		fmt.Println("AmfId From Models To NGAP Error: ", err.Error())
		return
	}
	shiftByte, err := GetBitString(tmpByte, 2, 6)
	if err != nil {
		fmt.Println("AmfId From Models To NGAP Error: ", err.Error())
		return
	}
	ptrId.BitLength = 6
	ptrId.Bytes = shiftByte
	return
}

func AmfIdToModels(regionId, setId, ptrId asn1rt.BitString) (amfId string) {
	regionHex := BitStringToHex(&regionId)
	tmpByte := []byte{setId.Bytes[0], (setId.Bytes[1] & 0xc0) | (ptrId.Bytes[0] >> 2)}
	restHex := hex.EncodeToString(tmpByte)
	amfId = regionHex + restHex
	return
}

// GetBitString is to get BitString with desire size from source byte array with bit offset
func GetBitString(srcBytes []byte, bitsOffset uint, numBits uint) (dstBytes []byte, err error) {
	bitsLeft := uint(len(srcBytes))*8 - bitsOffset
	if numBits > bitsLeft {
		err = fmt.Errorf("get bits overflow, requireBits: %d, leftBits: %d", numBits, bitsLeft)
		return
	}
	byteLen := (bitsOffset + numBits + 7) >> 3
	numBitsByteLen := (numBits + 7) >> 3
	dstBytes = make([]byte, numBitsByteLen)
	if numBitsByteLen == 0 {
		return
	}
	numBitsMask := byte(0xff)
	if modEight := numBits & 0x7; modEight != 0 {
		numBitsMask <<= uint8(8 - (modEight))
	}
	for i := 1; i < int(byteLen); i++ {
		dstBytes[i-1] = srcBytes[i-1]<<bitsOffset | srcBytes[i]>>(8-bitsOffset)
	}
	if byteLen == numBitsByteLen {
		dstBytes[byteLen-1] = srcBytes[byteLen-1] << bitsOffset
	}
	dstBytes[numBitsByteLen-1] &= numBitsMask
	return
}
