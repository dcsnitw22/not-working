package nas

import "errors"

func DecodePayloadContainerType(plcType byte) (string, error) {
	switch plcType {
	case N1SmInfo:
		return "N1 SM information", nil
	case Sms:
		return "SMS", nil
	case LppMsg:
		return "LTE Positioning Protocol message container", nil
	case SorTran:
		return "SOR transparent container", nil
	case UEPolCont:
		return "UE policy container", nil
	case UEParaUpdTran:
		return "UE parameters update transparent container", nil
	case LocaServMsg:
		return "Location services message container", nil
	case CIoTUser:
		return "CIoT user data container", nil
	case SLAACont:
		return "Service-level-AA container", nil
	case EventNotif:
		return "Event notification", nil
	case MultiPay:
		return "Multiple payloads", nil
	default:
		return "", errors.New("unknown PayLoad Container Type")
	}
}
