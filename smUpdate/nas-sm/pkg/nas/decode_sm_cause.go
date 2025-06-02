package nas

import "errors"

func DecodeSmCause(cause byte) (string, error) {
	switch cause {
	case InsufficientResources:
		return "INSUFFICIENT_RESOURCES", nil
	case MissingOrUnknownDNN:
		return "MISSING_OR_UNKNOWN_DNN", nil
	case UnknownPduSType:
		return "UNKNOWN_PDU_SESSION_TYPE", nil
	case UserAutheticationFail:
		return "USER_AUTHENTICATION_OR_AUTHORIZATION_FAILED", nil
	case RequestRejectedUnspecified:
		return "REQUEST_REJECTED_UNSPECIFIED", nil
	case ServiceOptionOutOfOrder:
		return "SERVICE_OPTION_TEMPORARILY_OUT_OF_ORDER", nil
	case PTIAlreadyInUse:
		return "PTI_ALREADY_IN_USE", nil
	case RegularDeactivation:
		return "REGULAR_DEACTIVATION", nil
	case NetworkFailure:
		return "NETWORK_FAILURE", nil
	case ReactivationRequested:
		return "REACTIVATION_REQUESTED", nil
	case SemanticErrorTFT:
		return "SEMANTIC_ERROR_IN_THE_TFT_OPERATION", nil
	case SyntacticalErrorTFT:
		return "SYNTACTICAL_ERROR_IN_THE_TFT_OPERATION", nil
	case InvalidPDUSIdentity:
		return "INVALID_PDU_SESSION_IDENTITY", nil
	case SemanticErrorPacketFilters:
		return "SEMANTIC_ERRORS_IN_PACKET_FILTERS", nil
	case SyntacticalErrorPacketFilters:
		return "SYNTACTICAL_ERROR_IN_PACKET_FILTERS", nil
	case OutOfLADNArea:
		return "OUT_OF_LADN_SERVICE_AREA", nil
	case PTIMismatch:
		return "PTI_MISMATCH", nil
	case PDUSTypeIPV4:
		return "PDU_SESSION_TYPE_IPV4_ONLY_ALLOWED", nil
	case PDUSTypeIPV6:
		return "PDU_SESSION_TYPE_IPV6_ONLY_ALLOWED", nil
	case PDUSDoesnotExist:
		return "PDU_SESSION_DOES_NOT_EXIST", nil
	case PDUSTypeIPV4V6:
		return "PDU_SESSION_TYPE_IPV4V6_ONLY_ALLOWED", nil
	case PDUSTypeUnstructured:
		return "PDU_SESSION_TYPE_UNSTRUCTURED_ONLY_ALLOWED", nil
	case Unsupported5QI:
		return "UNSUPPORTED_5QI_VALUE", nil
	case PDUSTypeEthernet:
		return "PDU_SESSION_TYPE_ETHERNET_ONLY_ALLOWED", nil
	case InsufficientResourcesSliceAndDNN:
		return "INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE_AND_DNN", nil
	case UnsupportedSSCMode:
		return "NOT_SUPPORTED_SSC_MODE", nil
	case InsufficientResourcesSlice:
		return "INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE", nil
	case MissingOrUnknownDNNSLice:
		return "MISSING_OR_UNKNOWN_DNN_IN_A_SLICE", nil
	case InvalidPTI:
		return "INVALID_PTI_VALUE", nil
	case MaxDataRateLow:
		return "MAXIMUM_DATA_RATE_PER_UE_FOR_USER_PLANE_INTEGRITY_PROTECTION_IS_TOO_LOW", nil
	case SemanticErrorQoS:
		return "SEMANTIC_ERROR_IN_THE_QOS_OPERATION", nil
	case SyntacticalErrorQoS:
		return "SYNTACTICAL_ERROR_IN_THE_QOS_OPERATION", nil
	case SemanticallyIncorrectMsg:
		return "SEMANTICALLY_INCORRECT_MESSAGE", nil
	case InvalidMandatoryInfo:
		return "INVALID_MANDATORY_INFORMATION", nil
	case MessageTypeNonExistent:
		return "MESSAGE_TYPE_NON_EXISTENT_OR_NOT_IMPLEMENTED", nil
	case MesageTypeNotCompatible:
		return "MESSAGE_TYPE_NOT_COMPATIBLE_WITH_THE_PROTOCOL_STATE", nil
	case IENonExistent:
		return "INFORMATION_ELEMENT_NON_EXISTENT_OR_NOT_IMPLEMENTED", nil
	case ConditionalIEError:
		return "CONDITIONAL_IE_ERROR", nil
	case MsgNotCompatibleProtocolState:
		return "MESSAGE_NOT_COMPATIBLE_WITH_THE_PROTOCOL_STATE", nil
	case ProtocolErrorUnspecified:
		return "PROTOCOL_ERROR_UNSPECIFIED", nil
	default:
		return "", errors.New("invalid SM Cause")
	}
}
