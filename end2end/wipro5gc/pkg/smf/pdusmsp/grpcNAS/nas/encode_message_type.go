package nas

import (
	"errors"
)

//Logic to encode MessageType parameter

func EncodeMessageType(messageType string) (byte, error) {
	switch messageType {
	case "REGISTRATION_REQUEST":
		return RegistrationRequest, nil
	case "REGISTRATION_ACCEPT":
		return RegistrationAccept, nil
	case "REGISTRATION_COMPLETE":
		return RegistrationComplete, nil
	case "REGISTRATION_REJECT":
		return RegistrationReject, nil
	case "DEREGISTRATION_REQUEST_UE_ORIGINATING":
		return DeregistrationRequestUeOrigin, nil
	case "DEREGISTRATION_ACCEPT_UE_ORIGINATING":
		return DeregistrationAcceptUeOrigin, nil
	case "DEREGISTRATION_REQUEST_UE_TERMINATED":
		return DeregistrationRequestUeTerminate, nil
	case "DEREGISTRATION_ACCEPT_UE_TERMINATED":
		return DeregistrationAcceptUeTerminate, nil
	case "SERVICE_REQUEST":
		return ServiceRequest, nil
	case "SERVICE_REJECT":
		return ServiceReject, nil
	case "SERVICE_ACCEPT":
		return ServiceAccept, nil
	case "CONFIGURATION_UPDATE_COMMAND":
		return ConfigurationUpdateCommand, nil
	case "CONFIGURATION_UPDATE_COMPLETE":
		return ConfigurationUpdateComplete, nil
	case "AUTHENTICATION_REQUEST":
		return AuthenticationRequest, nil
	case "AUTHENTICATION_RESPONSE":
		return AuthenticationResponse, nil
	case "AUTHENTICATION_REJECT":
		return AuthenticationReject, nil
	case "AUTHENTICATION_FAILURE":
		return AuthenticationFailure, nil
	case "AUTHENTICATION_RESULT":
		return AuthenticationResult, nil
	case "IDENTITY_REQUEST":
		return IdentityRequest, nil
	case "IDENTITY_RESPONSE":
		return IdentityResponse, nil
	case "SECURITY_MODE_COMMAND":
		return SecurtiyModeCommand, nil
	case "SECURITY_MODE_COMPLETE":
		return SecurtiyModeComplete, nil
	case "SECURITY_MODE_REJECT":
		return SecurtiyModeReject, nil
	case "FIVEG_MM_STATUS":
		return FiveGMMStatus, nil
	case "NOTIFICATION":
		return Notification, nil
	case "NOTIFICATION_RESPONSE":
		return NotificationResponse, nil
	case "UL_NAS_TRANSPORT":
		return UlNasTransport, nil
	case "DL_NAS_TRANSPORT":
		return DlNasTransport, nil
	case "PDU_SESSION_ESTABLISHMENT_REQUEST":
		return PduSEstablishmentRequest, nil
	case "PDU_SESSION_ESTABLISHMENT_ACCEPT":
		return PduSEstablishmentAccept, nil
	case "PDU_SESSION_ESTABLISHMENT_REJECT":
		return PduSEstablishmentReject, nil
	case "PDU_SESSION_AUTHENTICATION_COMMAND":
		return PduSAuthenticationCommand, nil
	case "PDU_SESSION_AUTHENTICATION_COMPLETE":
		return PduSAutheticationComplete, nil
	case "PDU_SESSION_AUTHENTICATION_RESULT":
		return PduSAuthenticationResult, nil
	case "PDU_SESSION_MODIFICATION_REQUEST":
		return PduSModificationRequest, nil
	case "PDU_SESSION_MODIFICATION_REJECT":
		return PduSModificationReject, nil
	case "PDU_SESSION_MODIFICATION_COMMAND":
		return PduSModificationCommand, nil
	case "PDU_SESSION_MODIFICATION_COMPLETE":
		return PduSModificationComplete, nil
	case "PDU_SESSION_MODIFICATION_COMMAND_REJECT":
		return PduSModificationCommandReject, nil
	case "PDU_SESSION_RELEASE_REQUEST":
		return PduSReleaseRequest, nil
	case "PDU_SESSION_RELEASE_REJECT":
		return PduSReleaseReject, nil
	case "PDU_SESSION_RELEASE_COMMAND":
		return PduSReleaseCommand, nil
	case "PDU_SESSION_RELEASE_COMPLETE":
		return PduSReleaseComplete, nil
	case "FIVEG_SM_STATUS":
		return FiveGSMStatus, nil
	default:
		return 0, errors.New("unknown MessageType")
	}
}
