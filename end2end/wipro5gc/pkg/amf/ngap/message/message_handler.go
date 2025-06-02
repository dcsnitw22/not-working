package message

import (
	"net"
	"os"

	"k8s.io/klog"

	// "github.com/free5gc/openapi/models"
	"w5gc.io/wipro5gcore/asn1gen"
	"w5gc.io/wipro5gcore/openapi/openapi_commn_client"
	"w5gc.io/wipro5gcore/pkg/amf/context"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/db/redis"
	"w5gc.io/wipro5gcore/pkg/amf/ngap/grpc"
)

//to do : initialise amfranpool
// var AmfRanPool [net.Conn]*AmfRan

var Ran *context.AmfRan

// routing function
func HandleMessage(conn net.Conn, msg []byte, grpc *grpc.Grpc, client *redis.RedisClient) {
	sdVal := "1"
	Ran = &context.AmfRan{
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
		klog.Info("NGAP Decode error: ", err.Error())
	} else {
		klog.Infoln("NGAPPDU unmarshaling done successfully")
	}
	klog.Infof("%+v\n", pdu)
	if pdu.T == asn1gen.NGAPPDUInitiatingMessageTAG {
		klog.Infof("%+v\n", *(pdu.U.InitiatingMessage))
		//NG Setup
		if pdu.U.InitiatingMessage.ProcedureCode == asn1gen.Asn1vIdNGSetup {
			klog.Info("RECEIVED NG SETUP MESSAGE")
			HandleNGSetupRequest(Ran, pdu)
		}
		//Uplink NAS
		if pdu.U.InitiatingMessage.ProcedureCode == asn1gen.Asn1vIdUplinkNASTransport {
			klog.Info("RECEIVED UPLINK NAS TRANSPORT MESSAGE")
			HandleUplinkNASTransport(Ran, pdu, grpc, client)
		}
		//Initial UE
		if pdu.U.InitiatingMessage.ProcedureCode == asn1gen.Asn1vIdInitialUEMessage {
			klog.Info("RECEIVED INITIAL UE MESSAGE")
			HandleInitialUEMessage(Ran, pdu, grpc, client)
		}
		//Downlink NAS
		/*if pdu.U.InitiatingMessage.ProcedureCode == asn1gen.Asn1vIdDownlinkNASTransport {
			klog.Info("RECEIVED DOWNLINK NAS TRANSPORT MESSAGE")
			HandleDownlinkNASTransport(ran,pdu)
		}*/
		//NG Reset initiated by NG RAN
		if pdu.U.InitiatingMessage.ProcedureCode == asn1gen.Asn1vIdNGReset {
			klog.Info("RECEIVED NG RESET MESSAGE")
			HandleNGResetMessage(Ran, pdu)
		}
		//NAS Non Delivery Indication
		//Error Indication
		if pdu.U.InitiatingMessage.ProcedureCode == asn1gen.Asn1vIdErrorIndication {
			klog.Info("RECEIVED ERROR INDICATION MESSAGE")
			HandleErrorIndicationFromRan(Ran, pdu)
		}
	} else if pdu.T == asn1gen.NGAPPDUSuccessfulOutcomeTAG {
		klog.Infof("%+v\n", *(pdu.U.SuccessfulOutcome))
		//NG Reset initiated by AMF; handle acknowledge
		if pdu.U.SuccessfulOutcome.ProcedureCode == asn1gen.Asn1vIdNGReset {
			klog.Info("RECEIVED NG RESET ACKNOWLEDGE")
			HandleNGResetAcknowledge(Ran, pdu)
		}

		//PDU Session Resource Setup Response
		if pdu.U.SuccessfulOutcome.ProcedureCode == asn1gen.Asn1vIdPDUSessionResourceSetup {
			klog.Info("RECEIVED PDU SESSION RESOURCE SETUP RESPONSE")
			HandlePduSessionResourceSetupResponse(Ran, pdu)
		}

		//Initial Context Setup Response
		if pdu.U.SuccessfulOutcome.ProcedureCode == asn1gen.Asn1vIdInitialContextSetup {
			klog.Info("RECEIVED INITIAL CONTEXT SETUP RESPONSE")
			os.WriteFile("initialContextSetupResponse", msg, 0644)
			HandleInitialContextSetupResponse(Ran, pdu)
		}
	} else if pdu.T == asn1gen.NGAPPDUUnsuccessfulOutcomeTAG {
		klog.Infof("%+v\n", *(pdu.U.UnsuccessfulOutcome))
		//Initial Context Setup Failure
		if pdu.U.SuccessfulOutcome.ProcedureCode == asn1gen.Asn1vIdInitialContextSetup {
			klog.Info("RECEIVED INITIAL CONTEXT SETUP FAILURE")
			HandleInitialContextSetupFailure(Ran, pdu)
		}
	} else {
		klog.Info("RECEIVED A MESSAGE")
	}
}
