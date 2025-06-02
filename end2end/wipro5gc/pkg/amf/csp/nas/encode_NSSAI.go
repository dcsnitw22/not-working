package nas

import "errors"

func EncodeNSSAI(nssai Nssai, byteArray []byte) ([]byte, error) {
	switch {
	case nssai.SST != nil && nssai.SD == nil && nssai.MappedHPLMNsd == nil && nssai.MappedHPLMNsst == nil:
		byteArray = append(byteArray, Sst)
		byteArray = append(byteArray, byte(*nssai.SST))
	case nssai.SST != nil && nssai.SD == nil && nssai.MappedHPLMNsd == nil && nssai.MappedHPLMNsst != nil:
		byteArray = append(byteArray, SstHPLMNSst)
		byteArray = append(byteArray, byte(*nssai.SST))
		byteArray = append(byteArray, byte(*nssai.MappedHPLMNsst))
	case nssai.SST != nil && nssai.SD != nil && nssai.MappedHPLMNsd == nil && nssai.MappedHPLMNsst == nil:
		byteArray = append(byteArray, SstSD)
		byteArray = append(byteArray, byte(*nssai.SST))

		//TODO: Add SD encoding here

	case nssai.SST != nil && nssai.SD != nil && nssai.MappedHPLMNsd == nil && nssai.MappedHPLMNsst != nil:
		byteArray = append(byteArray, SstSDHPLMNSst)
		byteArray = append(byteArray, byte(*nssai.SST))

		//TODO: Add SD encoding here

		byteArray = append(byteArray, byte(*nssai.MappedHPLMNsst))

	case nssai.SST != nil && nssai.SD != nil && nssai.MappedHPLMNsd != nil && nssai.MappedHPLMNsst != nil:
		byteArray = append(byteArray, All)
		byteArray = append(byteArray, byte(*nssai.SST))

		//TODO: Add SD encoding here
		byteArray = append(byteArray, byte(*nssai.MappedHPLMNsst))

		//TODO: Add HPLMNSD encoding here
	default:
		return nil, errors.New("invalid NSSAI")
	}

	return byteArray, nil

}
