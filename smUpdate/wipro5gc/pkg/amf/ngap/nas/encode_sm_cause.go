package nas

import (
	"errors"
)

//Logic to encode SM Cause

func EncodeSMcause(smCause string) (byte, error) {
	switch smCause {
	case "INSUFFICIENT_RESOURCES":
		return InsufficientResources, nil
	case "MISSING_OR_UNKNOWN_DNN":
		return MissingOrUnknownDNN, nil
	case "UNKNOWN_PDU_SESSION_TYPE":
		return UnknownPduSType, nil
	case "USER_AUTHENTICATION_OR_AUTHORIZATION_FAILED":
		return UserAutheticationFail, nil
	case "REQUEST_REJECTED_UNSPECIFIED":
		return RequestRejectedUnspecified, nil
	case "SERVICE_OPTION_TEMPORARILY_OUT_OF_ORDER":
		return ServiceOptionOutOfOrder, nil
	case "PTI_ALREADY_IN_USE":
		return PTIAlreadyInUse, nil
	case "REGULAR_DEACTIVATION":
		return RegularDeactivation, nil
	case "NETWORK_FAILURE":
		return NetworkFailure, nil
	case "REACTIVATION_REQUESTED":
		return ReactivationRequested, nil
	case "SEMANTIC_ERROR_IN_THE_TFT_OPERATION":
		return SemanticErrorTFT, nil
	case "SYNTACTICAL_ERROR_IN_THE_TFT_OPERATION":
		return SyntacticalErrorTFT, nil
	case "INVALID_PDU_SESSION_IDENTITY":
		return InvalidPDUSIdentity, nil
	case "SEMANTIC_ERRORS_IN_PACKET_FILTERS":
		return SemanticErrorPacketFilters, nil
	case "SYNTACTICAL_ERROR_IN_PACKET_FILTERS":
		return SyntacticalErrorPacketFilters, nil
	case "OUT_OF_LADN_SERVICE_AREA":
		return OutOfLADNArea, nil
	case "PTI_MISMATCH":
		return PTIMismatch, nil
	case "PDU_SESSION_TYPE_IPV4_ONLY_ALLOWED":
		return PDUSTypeIPV4, nil
	case "PDU_SESSION_TYPE_IPV6_ONLY_ALLOWED":
		return PDUSTypeIPV6, nil
	case "PDU_SESSION_DOES_NOT_EXIST":
		return PDUSDoesnotExist, nil
	case "PDU_SESSION_TYPE_IPV4V6_ONLY_ALLOWED":
		return PDUSTypeIPV4V6, nil
	case "PDU_SESSION_TYPE_UNSTRUCTURED_ONLY_ALLOWED":
		return PDUSTypeUnstructured, nil
	case "UNSUPPORTED_5QI_VALUE":
		return Unsupported5QI, nil
	case "PDU_SESSION_TYPE_ETHERNET_ONLY_ALLOWED":
		return PDUSTypeEthernet, nil
	case "INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE_AND_DNN":
		return InsufficientResourcesSliceAndDNN, nil
	case "NOT_SUPPORTED_SSC_MODE":
		return UnsupportedSSCMode, nil
	case "INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE":
		return InsufficientResourcesSlice, nil
	case "MISSING_OR_UNKNOWN_DNN_IN_A_SLICE":
		return MissingOrUnknownDNNSLice, nil
	case "INVALID_PTI_VALUE":
		return InvalidPTI, nil
	case "MAXIMUM_DATA_RATE_PER_UE_FOR_USER_PLANE_INTEGRITY_PROTECTION_IS_TOO_LOW":
		return MaxDataRateLow, nil
	case "SEMANTIC_ERROR_IN_THE_QOS_OPERATION":
		return SemanticErrorQoS, nil
	case "SYNTACTICAL_ERROR_IN_THE_QOS_OPERATION":
		return SyntacticalErrorQoS, nil
	case "SEMANTICALLY_INCORRECT_MESSAGE":
		return SemanticallyIncorrectMsg, nil
	case "INVALID_MANDATORY_INFORMATION":
		return InvalidMandatoryInfo, nil
	case "MESSAGE_TYPE_NON_EXISTENT_OR_NOT_IMPLEMENTED":
		return MessageTypeNonExistent, nil
	case "MESSAGE_TYPE_NOT_COMPATIBLE_WITH_THE_PROTOCOL_STATE":
		return MesageTypeNotCompatible, nil
	case "INFORMATION_ELEMENT_NON_EXISTENT_OR_NOT_IMPLEMENTED":
		return IENonExistent, nil
	case "CONDITIONAL_IE_ERROR":
		return ConditionalIEError, nil
	case "MESSAGE_NOT_COMPATIBLE_WITH_THE_PROTOCOL_STATE":
		return MsgNotCompatibleProtocolState, nil
	case "PROTOCOL_ERROR_UNSPECIFIED":
		return ProtocolErrorUnspecified, nil
	default:
		return 0, errors.New("invalid SM Cause")
	}
}
