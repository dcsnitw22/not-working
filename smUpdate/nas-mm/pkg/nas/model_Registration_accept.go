package nas

type RegistrationAcceptModel struct {
	ExtendedProtocolDiscriminator string
	SecurityHeaderType            string
	MessageType                   string
	RegResult                     string
	Sms                           string
	NssaPerformed                 string
	EmergencyReg                  string
	RoamingReg                    string
}
