package nas

type RegistrationRequestModel struct {
	ExtendedProtocolDiscriminator string
	SecurityHeaderType            string
	MessageType                   string
	RegistrationType              string
	FORValue                      string
	NASKSI                        string
	NASTSC                        string
	MobileIdentity                SupiIdentity
}
