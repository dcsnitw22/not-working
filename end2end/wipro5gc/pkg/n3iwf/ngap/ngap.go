package main

import (
	"errors"
	"net"

	"github.com/ishidawataru/sctp"
	"k8s.io/klog"
	"w5gc.io/wipro5gcore/asn1gen"
	"w5gc.io/wipro5gcore/pkg/n3iwf/ngap/message"
)

const rcvbuf int = 26144

func establishConnection(addr *sctp.SCTPAddr) error {
	conn, err := sctp.DialSCTP("sctp", nil, addr)
	if err != nil {
		return errors.New("failed to dial : " + err.Error())
	}
	klog.Infof("Connected to RemoteAddr: %s", conn.RemoteAddr())
	/*wconn := sctp.NewSCTPSndRcvInfoWrappedConn(conn)
	if rcvbuf != 0 {
		err = wconn.SetReadBuffer(rcvbuf)
		if err != nil {
			klog.Fatalf("failed to set read buf: %v", err)
			return err
		}
	}*/
	// Send N2 Setup Request after connection is established
	//N2 Setup = NG Setup
	err = message.SendN2SetupRequest(conn)
	if err != nil {
		return errors.New("failed to send N2 Setup Request : " + err.Error())
	}
	go handleConnection(conn)
	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, rcvbuf+128)
		_, err := conn.Read(buf)
		if err != nil {
			klog.Errorf("Read failed : %v", err)
			return
		}
		//handle message now
		handleMessage(conn, buf)
	}
}

func handleMessage(conn net.Conn, buf []byte) {
	ngapPDU := &asn1gen.NGAPPDU{}
	_, err := asn1gen.Unmarshal(buf, ngapPDU)
	if err != nil {
		klog.Error("failed to decode NGAP message : ", err)
		return
	}
	if ngapPDU == nil {
		klog.Error("NGAP message is nil")
		return
	}
	klog.Info("Received message: %v", ngapPDU)
	switch ngapPDU.T {
	case asn1gen.NGAPPDUInitiatingMessageTAG:
		switch ngapPDU.U.InitiatingMessage.ProcedureCode {
		case asn1gen.Asn1vIdInitialContextSetup:
			klog.Info("Received Initial Context Setup Request")
			message.HandleInitialContextSetupRequest(conn, ngapPDU)
		case asn1gen.Asn1vIdDownlinkNASTransport:
			klog.Info("Received Downlink NAS Transport")
			message.HandleDownlinkNASTransport(conn, ngapPDU)
		case asn1gen.Asn1vIdPDUSessionResourceSetup:
			klog.Info("Received PDU Session Resource Setup Request")
			message.HandlePduSessionResourceSetupRequest(conn, ngapPDU)
		default:
			klog.Info("Received unknown message")
		}
	case asn1gen.NGAPPDUSuccessfulOutcomeTAG:
		switch ngapPDU.U.SuccessfulOutcome.ProcedureCode {
		case asn1gen.Asn1vIdNGSetup:
			klog.Info("Received NG Setup Response")
			message.HandleNGSetupResponse(conn, ngapPDU)
		default:
			klog.Info("Received unknown message")
		}
	case asn1gen.NGAPPDUUnsuccessfulOutcomeTAG:
		switch ngapPDU.U.UnsuccessfulOutcome.ProcedureCode {
		case asn1gen.Asn1vIdNGSetup:
			klog.Info("Received NG Setup Failure")
			message.HandleNGSetupFailure(conn, ngapPDU)
		default:
			klog.Info("Received unknown message")
		}
	default:
		klog.Info("Received unknown message")
	}
}

func main() {
	addr := &sctp.SCTPAddr{
		IPAddrs: []net.IPAddr{
			{IP: net.ParseIP("10.100.10.230")},
		},
		Port: 38412,
	}
	err := establishConnection(addr)
	if err != nil {
		klog.Fatalf("failed to establish connection: %v", err)
	}
}
