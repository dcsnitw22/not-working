package main

import (
	"context"
	"fmt"
	"log"
	"nasMM/pkg/ngapNas/pb"
	// "os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/anypb"
)

func Trail() {
	serverAddr := "grpcnasamf-service.nasmm.svc.cluster.local:50052"

	clientConn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer clientConn.Close()

	client := pb.NewDataServiceClient(clientConn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// registrationRequestFile, _ := os.OpenFile("/home/wipro/mounika-trail/mounika_nas/nas-sm/testFiles/MinRegistrationRequest", os.O_RDONLY, 0)
	// fileInfo, _ := registrationRequestFile.Stat()
	// registrationRequestByteArray := make([]byte, fileInfo.Size())
	// _, _ = registrationRequestFile.Read(registrationRequestByteArray)
	// fmt.Println("Byte Array:", registrationRequestByteArray)
	val := []byte{126, 0, 65, 121, 0, 13, 1, 130, 246, 16, 0, 0, 0, 0, 0, 0, 0, 0, 16, 46, 4, 240, 240, 240, 240}
    l := 25
    reg := make([]byte, l)
    copy(reg, val)

	ba := pb.NasMessage{NasPdu: val}
	fmt.Println(ba)

	anyMessage, err := anypb.New(&ba)
	if err != nil {
		log.Fatalf("Failed to create Any message: %v", err)
	}

	fmt.Println(anyMessage)

	data := pb.DataRequest{Data: anyMessage, ReqType: "Registration Request"}

	fmt.Println(&data)

	resp, err := client.SendData(ctx, &data)

	if err != nil {
		log.Fatalf("Error sending SMData: %v", err)
	}

	fmt.Println(resp)
}

func main() {
	Trail()
}
