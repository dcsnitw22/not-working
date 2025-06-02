package nas

import (
	"errors"
)

//Logic to decode MessageType parameter

func DecodeMessageType(messageType byte) (string, error) {
	switch messageType {
	case RegistrationRequest:
		return "REGISTRATION_REQUEST", nil
	case RegistrationAccept:
		return "REGISTRATION_ACCEPT", nil
	case RegistrationComplete:
		return "REGISTRATION_COMPLETE", nil
	case RegistrationReject:
		return "REGISTRATION_REJECT", nil
	case DeregistrationRequestUeOrigin:
		return "DEREGISTRATION_REQUEST_UE_ORIGINATING", nil
	case DeregistrationAcceptUeOrigin:
		return "DEREGISTRATION_ACCEPT_UE_ORIGINATING", nil
	case DeregistrationRequestUeTerminate:
		return "DEREGISTRATION_REQUEST_UE_TERMINATED", nil
	case DeregistrationAcceptUeTerminate:
		return "DEREGISTRATION_ACCEPT_UE_TERMINATED", nil
	case ServiceRequest:
		return "SERVICE_REQUEST", nil
	case ServiceReject:
		return "SERVICE_REJECT", nil
	case ServiceAccept:
		return "SERVICE_ACCEPT", nil
	case ConfigurationUpdateCommand:
		return "CONFIGURATION_UPDATE_COMMAND", nil
	case ConfigurationUpdateComplete:
		return "CONFIGURATION_UPDATE_COMPLETE", nil
	case AuthenticationRequest:
		return "AUTHENTICATION_REQUEST", nil
	case AuthenticationResponse:
		return "AUTHENTICATION_RESPONSE", nil
	case AuthenticationReject:
		return "AUTHENTICATION_REJECT", nil
	case AuthenticationFailure:
		return "AUTHENTICATION_FAILURE", nil
	case AuthenticationResult:
		return "AUTHENTICATION_RESULT", nil
	case IdentityRequest:
		return "IDENTITY_REQUEST", nil
	case IdentityResponse:
		return "IDENTITY_RESPONSE", nil
	case SecurtiyModeCommand:
		return "SECURITY_MODE_COMMAND", nil
	case SecurtiyModeComplete:
		return "SECURITY_MODE_COMPLETE", nil
	case SecurtiyModeReject:
		return "SECURITY_MODE_REJECT", nil
	case FiveGMMStatus:
		return "FIVEG_MM_STATUS", nil
	case Notification:
		return "NOTIFICATION", nil
	case NotificationResponse:
		return "NOTIFICATION_RESPONSE", nil
	case UlNasTransport:
		return "UL_NAS_TRANSPORT", nil
	case DlNasTransport:
		return "DL_NAS_TRANSPORT", nil
	case PduSEstablishmentRequest:
		return "PDU_SESSION_ESTABLISHMENT_REQUEST", nil
	case PduSEstablishmentAccept:
		return "PDU_SESSION_ESTABLISHMENT_ACCEPT", nil
	case PduSEstablishmentReject:
		return "PDU_SESSION_ESTABLISHMENT_REJECT", nil
	case PduSAuthenticationCommand:
		return "PDU_SESSION_AUTHENTICATION_COMMAND", nil
	case PduSAutheticationComplete:
		return "PDU_SESSION_AUTHENTICATION_COMPLETE", nil
	case PduSAuthenticationResult:
		return "PDU_SESSION_AUTHENTICATION_RESULT", nil
	case PduSModificationRequest:
		return "PDU_SESSION_MODIFICATION_REQUEST", nil
	case PduSModificationReject:
		return "PDU_SESSION_MODIFICATION_REJECT", nil
	case PduSModificationCommand:
		return "PDU_SESSION_MODIFICATION_COMMAND", nil
	case PduSModificationComplete:
		return "PDU_SESSION_MODIFICATION_COMPLETE", nil
	case PduSModificationCommandReject:
		return "PDU_SESSION_MODIFICATION_COMMAND_REJECT", nil
	case PduSReleaseRequest:
		return "PDU_SESSION_RELEASE_REQUEST", nil
	case PduSReleaseReject:
		return "PDU_SESSION_RELEASE_REJECT", nil
	case PduSReleaseCommand:
		return "PDU_SESSION_RELEASE_COMMAND", nil
	case PduSReleaseComplete:
		return "PDU_SESSION_RELEASE_COMPLETE", nil
	case FiveGSMStatus:
		return "FIVEG_SM_STATUS", nil
	default:
		return "", errors.New("unknown MessageType")
	}
}
