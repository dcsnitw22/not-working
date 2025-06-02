package nas

import "errors"

func EncodePacketFilterComponentTypeIdentifier(pfci string) (byte, error) {
	switch pfci {
	case "Match-all type":
		return MatchAllType, nil
	case "IPv4 remote address type":
		return IPv4RemoteAddressType, nil
	case "IPv4 local address type":
		return IPv4LocalAddressType, nil
	case "IPv6 remote address/prefix length type":
		return IPv6RemoteAddress, nil
	case "IPv6 local address/prefix length type":
		return IPv6LocalAddress, nil
	case "Protocol identifier/Next header type":
		return ProtocolIdentifier, nil
	case "Single local port type":
		return SingleLocalPortType, nil
	case "Local port range type":
		return LocalPortRangeType, nil
	case "Single remote port type":
		return SingleRemotePortType, nil
	case "Remote port range type":
		return RemotePortRangeType, nil
	case "Security parameter index type":
		return SecurityParameterIndexType, nil
	case "Type of service/Traffic class type":
		return TypeOfService, nil
	case "Flow label type":
		return FlowLabelType, nil
	case "Destination MAC address type":
		return DestinationMACaddressType, nil
	case "Source MAC address type":
		return SourceMACAddressType, nil
	case "802.1Q C-TAG VID type":
		return CTAGVIDtype, nil
	case "802.1Q S-TAG VID type":
		return STAGVIDtype, nil
	case "802.1Q C-TAG PCP/DEI type":
		return CTAGPCPtype, nil
	case "802.1Q S-TAG PCP/DEI type":
		return STAGPCPtype, nil
	case "Ethertype type":
		return Ethertype, nil
	case "Destination MAC address range type":
		return DestinationMAC, nil
	case "Source MAC address range type":
		return SourceMAC, nil
	}
	return 0, errors.New("invalid Input")
}
