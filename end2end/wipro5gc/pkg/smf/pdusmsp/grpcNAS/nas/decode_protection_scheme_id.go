package nas

func DecodeProtectionSchemeID(id byte)(string){
	switch id{
	case NullScheme:
		return "Null scheme"
	case ECIESA:
		return "ECIES scheme profile A"
	case ECISESB:
		return "ECIES scheme profile B"
	default:
		return "Reserved"
	}
}
