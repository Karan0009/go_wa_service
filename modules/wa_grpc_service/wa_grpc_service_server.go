package wa_grpc_service

import (
	"context"
	"fmt"

	wa_grpc_service "github.com/Karan0009/go_wa_bot/modules/wa_grpc_service/proto"
	"github.com/Karan0009/go_wa_bot/modules/wa_service"
)

type MyWaGrpcServiceServer struct {
	wa_grpc_service.UnimplementedWAServiceServer
}

func GetMyWaGrpcServiceServer() *MyWaGrpcServiceServer {
	return &MyWaGrpcServiceServer{}
}

func (s *MyWaGrpcServiceServer) SendOtpMessage(
	context.Context,
	*wa_grpc_service.SendOtpRequest) (*wa_grpc_service.SendOtpResponse, error) {
	waClient, waClientErr := wa_service.GetWaClient()
	if waClientErr != nil {
		return &wa_grpc_service.SendOtpResponse{StatusCode: "500", Success: false, Data: fmt.Sprintf("%v", waClientErr)}, nil
	}
	err := waClient.SendMessage("917015295819", "boooooo")
	if err != nil {
		return &wa_grpc_service.SendOtpResponse{StatusCode: "500", Success: false, Data: fmt.Sprintf("%v", err)}, nil
	}
	return &wa_grpc_service.SendOtpResponse{StatusCode: "200", Success: true, Data: ""}, nil
}
