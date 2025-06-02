package nas

import "errors"

func EncodePayloadContainerType(plcType string) (byte, error) {
	switch plcType {
	case "N1 SM information":
		return N1SmInfo, nil
	case "SMS":
		return Sms, nil
	case "LTE Positioning Protocol message container":
		return LppMsg, nil
	case "SOR transparent container":
		return SorTran, nil
	case "UE policy container":
		return UEPolCont, nil
	case "UE parameters update transparent container":
		return UEParaUpdTran, nil
	case "Location services message container":
		return LocaServMsg, nil
	case "CIoT user data container":
		return CIoTUser, nil
	case "Service-level-AA container":
		return SLAACont, nil
	case "Event notification":
		return EventNotif, nil
	case "Multiple payloads":
		return MultiPay, nil
	default:
		return 0, errors.New("invalid Payload Container Type")
	}

}
