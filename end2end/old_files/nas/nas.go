package nas

import (
	"fmt"
	"os"
)

func DecodeInitialUE(registrationRequestFile *os.File) error {
	fmt.Println("Decode Registration Request")
	// var b = make([]byte, 25)
	// _, err := registrationFile.Read(b)
	// registrationFile.Seek(0, 0)
	// fmt.Println(b, n, err)

	// Get the file size
	fileInfo, _ := registrationRequestFile.Stat()
	fileSize := fileInfo.Size()
	// Create a byte array with the size of the file
	registrationRequestByteArray := make([]byte, fileSize)
	// Read the bytes from the file into the byte array
	_, _ = registrationRequestFile.Read(registrationRequestByteArray)
	fmt.Println("Byte Array:", registrationRequestByteArray)
	registrationPDU, err := DecodeRegistrationRequest(registrationRequestByteArray)
	fmt.Printf("Decoded PDU Session Registration Request: %+v\n", registrationPDU)
	fmt.Println("Error while decoding PDU Session Registration Request:", err)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func EncodeRegistrationAcceptMain() *os.File {
	fmt.Println("Encode Registration Accept")
	registrationAccept := RegistrationAcceptModel{ExtendedProtocolDiscriminator: "MOBILITY_MANAGEMENT_MESSAGES", SecurityHeaderType: "Plain 5GS NAS message", MessageType: "REGISTRATION_ACCEPT", RegResult: "3GPP access", Sms: "SMS over NAS not allowed", NssaPerformed: "Network slice-specific authentication and authorization is not to be performed", EmergencyReg: "Not registered for emergency services", RoamingReg: "No additional information"}
	//file, error := EncodeRegistrationAccept("binaryRegistrationAccept", regAcc)
	registrationAcceptByteArray, err := EncodeRegistrationAccept(registrationAccept)
	fmt.Println("Encoded Registartion Accept: ", registrationAcceptByteArray)
	fmt.Println("Error while encoding Registration accept:", err)
	file, err := os.Create("binaryRegistrationAccept")
	if err != nil {
		fmt.Println("Error while creating Registration accept file:", err)
		return nil
	}
	n, err := file.Write(registrationAcceptByteArray)
	if err != nil {
		fmt.Println("Error while writing to Registration accept file:", err)
		return nil
	}
	fmt.Printf("wrote %d bytes\n", n)
	return file
}

/*func main() {

	fmt.Println("### Decode Mobility Management Procedures ###")

	fmt.Println()

	//Registration Request
	fmt.Println("### Decode Registration Request ###")
	registrationRequestFile, _ := os.OpenFile("/home/ubuntu/Desktop/NAS/testFiles/MinRegistrationRequest", os.O_RDONLY, 0)
	// Get the file size
	fileInfo, _ := registrationRequestFile.Stat()
	fileSize := fileInfo.Size()
	// Create a byte array with the size of the file
	registrationRequestByteArray := make([]byte, fileSize)
	// Read the bytes from the file into the byte array
	_, _ = registrationRequestFile.Read(registrationRequestByteArray)
	fmt.Println("Byte Array:", registrationRequestByteArray)
	registrationPDU, err := nas.DecodeRegistrationRequest(registrationRequestByteArray)
	fmt.Printf("Decoded PDU Session Registration Request: %+v\n", registrationPDU)
	fmt.Println("Error while decoding PDU Session Registration Request:", err)

	fmt.Println()

	//UpLink NAS
	fmt.Println("### Decode Uplink NAS message ###")
	ulNasFile, _ := os.OpenFile("/home/ubuntu/Desktop/NAS/testFiles/MinULNASTransport", os.O_RDONLY, 0)
	// Get the file size
	fileInfo, _ = ulNasFile.Stat()
	fileSize = fileInfo.Size()
	// Create a byte array with the size of the file
	ulNasByteArray := make([]byte, fileSize)
	// Read the bytes from the file into the byte array
	_, _ = ulNasFile.Read(ulNasByteArray)
	fmt.Println("Byte Array:", ulNasByteArray)
	ulNasPDU, err := nas.DecodeULNas(ulNasByteArray)
	fmt.Printf("Decoded UL Nas Message:  %+v\n", ulNasPDU)
	fmt.Println("Error while decoding UL Nas Message:", err)

	fmt.Println()

	fmt.Println("### Encoding Mobility Management Procedures ###")

	fmt.Println()

	fmt.Println("### Encode Registration Accept ###")
	registrationAccept := nas.RegistrationAcceptModel{ExtendedProtocolDiscriminator: "MOBILITY_MANAGEMENT_MESSAGES", SecurityHeaderType: "Plain 5GS NAS message", MessageType: "REGISTRATION_ACCEPT", RegResult: "3GPP access", Sms: "SMS over NAS not allowed", NssaPerformed: "Network slice-specific authentication and authorization is not to be performed", EmergencyReg: "Not registered for emergency services", RoamingReg: "No additional information"}
	fmt.Printf("Dummy Data for Registartion Accept: %+v\n", registrationAccept)
	registrationAcceptByteArray, err := nas.EncodeRegistrationAccept(registrationAccept)
	fmt.Println("Encoded Registartion Accept: ", registrationAcceptByteArray)
	fmt.Println("Error while encoding Registration accept:", err)

	fmt.Println()

	fmt.Println("### Encode DL NAS Message ###")
	dlNAS := nas.DLNasModel{ExtendedProtocolDiscriminator: "MOBILITY_MANAGEMENT_MESSAGES", SecurityHeaderType: "Plain 5GS NAS message", MessageType: "DL_NAS_TRANSPORT", PayLoadContainerType: "N1 SM information"}
	fmt.Printf("Dummy Data for DL NAS Transport: %+v\n", dlNAS)
	dlNasByteArray, err := nas.EncodeDLNAS(dlNAS)
	fmt.Println("Encoded DL NAS Transport:", dlNasByteArray)
	fmt.Println("Error while encoding DL NAS Transport:", err)

	fmt.Println()

	fmt.Println("### Decode Session Management Procedures ###")

	fmt.Println()

	//PDU session establishment request
	fmt.Println("### Decode PDU Session Establishment Request ###")
	establishmentRequestFile, _ := os.OpenFile("/home/ubuntu/Desktop/NAS/testFiles/MinPDUSessionEstablishmentRequest", os.O_RDONLY, 0)
	// Get the file size
	fileInfo, _ = establishmentRequestFile.Stat()
	fileSize = fileInfo.Size()
	// Create a byte array with the size of the file
	establishmentRequestByteArray := make([]byte, fileSize)
	// Read the bytes from the file into the byte array
	_, _ = establishmentRequestFile.Read(establishmentRequestByteArray)
	fmt.Println("Byte Array:", establishmentRequestByteArray)
	establishmentPDU, err := nas.DecodePduSessionEstablishmentRequest(establishmentRequestByteArray)
	fmt.Printf("Decoded PDU Session Establishment Request: %+v\n", establishmentPDU)
	fmt.Println("Error while decoding PDU Session Establishment Request:", err)

	fmt.Println()

	//PDU Session Modification Request
	fmt.Println("### Decode PDU Session Modification Request ###")
	modificationRequestFile, _ := os.OpenFile("/home/ubuntu/Desktop/NAS/testFiles/MinPDUSessionModificationRequest", os.O_RDONLY, 0)
	// Get the file size
	fileInfo, _ = modificationRequestFile.Stat()
	fileSize = fileInfo.Size()
	// Create a byte array with the size of the file
	modificationRequestByteArray := make([]byte, fileSize)
	// Read the bytes from the file into the byte array
	_, _ = modificationRequestFile.Read(modificationRequestByteArray)
	fmt.Println("Byte Array:", modificationRequestByteArray)
	modificationPDU, err := nas.DecodePduSessionModificationRequest(modificationRequestByteArray)
	fmt.Printf("Decoded PDU Session Modification Request: %+v\n", modificationPDU)
	fmt.Println("Error while decoding PDU Session Modification Request:", err)

	fmt.Println()

	//PDU Session Release Request
	fmt.Println("### Decode PDU Session Release Request ###")
	releaseRequestFile, _ := os.OpenFile("/home/ubuntu/Desktop/NAS/testFiles/MinPDUSessionReleaseRequest", os.O_RDONLY, 0)
	// Get the file size
	fileInfo, _ = releaseRequestFile.Stat()
	fileSize = fileInfo.Size()
	// Create a byte array with the size of the file
	releaseRequestByteArray := make([]byte, fileSize)
	// Read the bytes from the file into the byte array
	_, _ = releaseRequestFile.Read(releaseRequestByteArray)
	fmt.Println("Byte Array:", releaseRequestByteArray)
	releasePDU, err := nas.DecodePduSessionReleaseRequest(releaseRequestByteArray)
	fmt.Printf("Decoded PDU Session Release Request: %+v\n", releasePDU)
	fmt.Println("Error while decoding PDU Session Release Request:", err)

	fmt.Println()

	fmt.Println("### Encoding Session Management Procedures")
	fmt.Println()

	//PDU Session Establishment Accept
	fmt.Println("### Encode PDU Session Establishment Accept ###")
	establishmentAccept := nas.PduSessionEstablishmentAccept{ExtendedProtocolDiscriminator: "SESSION_MANAGEMENT_MESSAGES", PDUsessionId: "VAL_1", PTI: 4, MessageType: "PDU_SESSION_ESTABLISHMENT_ACCEPT", PduSessionType: "IPV4", SSCmode: "SSC_MODE_1", SessionAmbr: "MULT_1Kbps"}
	fmt.Printf("Dummy Data for Establishment Accept: %+v\n", establishmentAccept)
	establishmentAcceptByteArray, err := nas.EncodePduSessionEstablishmentAccept(establishmentAccept)
	fmt.Println("Encoded Session Establishment Accept:", establishmentAcceptByteArray)
	fmt.Println("Error while encoding PDU Session Establishment Accept:", err)

	// dummyestablishmentAcceptFile, _ := os.OpenFile("/home/ubuntu/Desktop/NAS/testFiles/MinPDUSessionEstablishmentAccept", os.O_RDONLY, 0)
	// // Get the file size
	// fileInfo, _ = dummyestablishmentAcceptFile.Stat()
	// fileSize = fileInfo.Size()
	// // Create a byte array with the size of the file
	// dummyestablishmentAcceptByteArray := make([]byte, fileSize)
	// // Read the bytes from the file into the byte array
	// _, _ = dummyestablishmentAcceptFile.Read(dummyestablishmentAcceptByteArray)
	// fmt.Println("Byte Array:", dummyestablishmentAcceptByteArray)

	fmt.Println()

	//PDU Session Establishment Reject
	fmt.Println("### Encode PDU Session Establishment Reject ###")
	establishmentReject := nas.PduSessionEstablishmentReject{ExtendedProtocolDiscriminator: "SESSION_MANAGEMENT_MESSAGES", PDUsessionId: "VAL_1", PTI: 4, MessageType: "PDU_SESSION_ESTABLISHMENT_REJECT", SMcause: "INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE_AND_DNN"}
	fmt.Printf("Dummy Data for Establishment Reject: %+v\n", establishmentReject)
	establishmentRejectByteArray, err := nas.EncodePduSessionEstablishmentReject(establishmentReject)
	fmt.Println("Encoded Session Establishment Reject:", establishmentRejectByteArray)
	fmt.Println("Error while encoding PDU Session Establishment Reject:", err)

	fmt.Println()

	//PDU Session Modification Reject
	fmt.Println("### Encode PDU Session Modification Reject ###")
	modificationReject := nas.PduSessionModificationReject{ExtendedProtocolDiscriminator: "SESSION_MANAGEMENT_MESSAGES", PDUsessionId: "VAL_1", PTI: 4, MessageType: "PDU_SESSION_MODIFICATION_REJECT", SMcause: "INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE_AND_DNN"}
	fmt.Printf("Dummy Data for Modification Reject: %+v\n", modificationReject)
	modificationRejectByteArray, err := nas.EncodePduSessionModificationReject(modificationReject)
	fmt.Println("Encoded Session Modification Reject:", modificationRejectByteArray)
	fmt.Println("Error while encoding PDU Session Modification Reject:", err)

	fmt.Println()

	//PDU Session Release Reject
	fmt.Println("### Encode PDU Session Release Reject ###")
	releaseReject := nas.PduSessionReleaseReject{ExtendedProtocolDiscriminator: "SESSION_MANAGEMENT_MESSAGES", PDUsessionId: "VAL_1", PTI: 4, MessageType: "PDU_SESSION_RELEASE_REJECT", SMcause: "INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE_AND_DNN"}
	fmt.Printf("Dummy Data for Release Reject: %+v\n", releaseReject)
	releaseRejectByteArray, err := nas.EncodePduSessionReleaseReject(releaseReject)
	fmt.Println("Encoded Session Release Reject:", releaseRejectByteArray)
	fmt.Println("Error while encoding PDU Session Release Reject:", err)

	fmt.Println()

}*/
