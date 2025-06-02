package message

import (
	"net"
	"os"

	// "github.com/free5gc/openapi/models"
	"w5gc.io/wipro5gcore/openapi/openapi_commn_client"
	"w5gc.io/wipro5gcore/pkg/amf/context"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/asn1gen"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/grpc"
	logger "w5gc.io/wipro5gcore/pkg/amf/ngap/log"
)

//to do : initialise amfranpool
// var AmfRanPool [net.Conn]*AmfRan

// to do : error handling
// check whether unmarshaling is correct
func HandleMessage(conn net.Conn, msg []byte, grpc *grpc.Grpc) {
	sdVal := "1"
	ran := &context.AmfRan{
		RanId: &openapi_commn_client.GlobalRanNodeId{
			PlmnId: &openapi_commn_client.PlmnId{
				Mcc: "286",
				Mnc: "01",
			},
			GNbId: &openapi_commn_client.GNbId{
				BitLength: 24,
				GNBValue:  "000102",
			},
		},

		Name: "ueransim",

		AnType: "3GPP_ACCESS",

		Conn: conn,

		SupportedTAList: []context.SupportedTAI{
			{
				Tai: openapi_commn_client.Tai{
					PlmnId: openapi_commn_client.PlmnId{
						Mcc: "286",
						Mnc: "01",
					},
					Tac: "000001",
				},
				SNssaiList: []openapi_commn_client.Snssai{
					{
						Sst: 1,
						Sd:  &sdVal,
					},
				},
			},
		},
	}
	// var pdu *asn1gen.NGAPPDU won't work
	// var has to be initialised
	var pdu *asn1gen.NGAPPDU = &asn1gen.NGAPPDU{}
	_, err := asn1gen.Unmarshal(msg, pdu)
	if err != nil {
		println("NGAP Decode error: ", err.Error())
	} else {
		logger.AppLog.Infoln("NGAPPDU unmarshaling done successfully")
	}
	logger.AppLog.Infof("%+v\n", pdu)
	if pdu.T == asn1gen.NGAPPDUInitiatingMessageTAG {
		logger.AppLog.Infof("%+v\n", *(pdu.U.InitiatingMessage))
		//NG Setup
		if pdu.U.InitiatingMessage.ProcedureCode == asn1gen.Asn1vIdNGSetup {
			logger.AppLog.Println("RECEIVED NG SETUP MESSAGE")
			HandleNGSetupRequest(ran, pdu)
		}
		//Uplink NAS
		if pdu.U.InitiatingMessage.ProcedureCode == asn1gen.Asn1vIdUplinkNASTransport {
			logger.AppLog.Println("RECEIVED UPLINK NAS TRANSPORT MESSAGE")
			HandleUplinkNASTransport(ran, pdu, grpc)
		}
		//Initial UE
		if pdu.U.InitiatingMessage.ProcedureCode == asn1gen.Asn1vIdInitialUEMessage {
			logger.AppLog.Println("RECEIVED INITIAL UE MESSAGE")
			HandleInitialUEMessage(ran, pdu, grpc)
		}
		//Downlink NAS
		/*if pdu.U.InitiatingMessage.ProcedureCode == asn1gen.Asn1vIdDownlinkNASTransport {
			logger.AppLog.Println("RECEIVED DOWNLINK NAS TRANSPORT MESSAGE")
			HandleDownlinkNASTransport(ran,pdu)
		}*/
		//NG Reset initiated by NG RAN
		if pdu.U.InitiatingMessage.ProcedureCode == asn1gen.Asn1vIdNGReset {
			logger.AppLog.Println("RECEIVED NG RESET MESSAGE")
			HandleNGResetMessage(ran, pdu)
		}
		//NAS Non Delivery Indication
		//Error Indication
		if pdu.U.InitiatingMessage.ProcedureCode == asn1gen.Asn1vIdErrorIndication {
			logger.AppLog.Println("RECEIVED ERROR INDICATION MESSAGE")
			HandleErrorIndicationFromRan(ran, pdu)
		}
	} else if pdu.T == asn1gen.NGAPPDUSuccessfulOutcomeTAG {
		logger.AppLog.Infof("%+v\n", *(pdu.U.SuccessfulOutcome))
		//NG Reset initiated by AMF; handle acknowledge
		if pdu.U.SuccessfulOutcome.ProcedureCode == asn1gen.Asn1vIdNGReset {
			logger.AppLog.Println("RECEIVED NG RESET ACKNOWLEDGE")
			HandleNGResetAcknowledge(ran, pdu)
		}

		//PDU Session Resource Setup Response
		if pdu.U.SuccessfulOutcome.ProcedureCode == asn1gen.Asn1vIdPDUSessionResourceSetup {
			logger.AppLog.Println("RECEIVED PDU SESSION RESOURCE SETUP RESPONSE")
			HandlePduSessionResourceSetupResponse(ran, pdu)
		}

		//Initial Context Setup Response
		if pdu.U.SuccessfulOutcome.ProcedureCode == asn1gen.Asn1vIdInitialContextSetup {
			logger.AppLog.Println("RECEIVED INITIAL CONTEXT SETUP RESPONSE")
			os.WriteFile("initialContextSetupResponse", msg, 0644)
			HandleInitialContextSetupResponse(ran, pdu)
		}
	} else if pdu.T == asn1gen.NGAPPDUUnsuccessfulOutcomeTAG {
		logger.AppLog.Infof("%+v\n", *(pdu.U.UnsuccessfulOutcome))
		//Initial Context Setup Failure
		if pdu.U.SuccessfulOutcome.ProcedureCode == asn1gen.Asn1vIdInitialContextSetup {
			logger.AppLog.Println("RECEIVED INITIAL CONTEXT SETUP FAILURE")
			HandleInitialContextSetupFailure(ran, pdu)
		}
	} else {
		logger.AppLog.Println("RECEIVED A MESSAGE")
	}
}
