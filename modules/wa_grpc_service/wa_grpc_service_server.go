package wa_grpc_service

import (
	"context"
	"fmt"

	wa_grpc_service "github.com/Karan0009/go_wa_bot/modules/wa_grpc_service/proto"
	"github.com/Karan0009/go_wa_bot/modules/wa_service"
	"google.golang.org/grpc/codes"
)

type MyWaGrpcServiceServer struct {
	wa_grpc_service.UnimplementedWAServiceServer
}

func GetMyWaGrpcServiceServer() *MyWaGrpcServiceServer {
	return &MyWaGrpcServiceServer{}
}

func (s *MyWaGrpcServiceServer) SendOtpMessage(
	ctx context.Context,
	req *wa_grpc_service.SendOtpRequest) (*wa_grpc_service.SendOtpResponse, error) {
	otpMessage := *req.GetData()
	if len(otpMessage.PhoneNumber) != 12 {
		return &wa_grpc_service.SendOtpResponse{StatusCode: codes.InvalidArgument.String(), Success: false, Data: fmt.Sprintf("%v", "phone must be prefixed with 91")}, nil
	}
	if len(otpMessage.OtpCode) == 0 {
		return &wa_grpc_service.SendOtpResponse{StatusCode: codes.InvalidArgument.String(), Success: false, Data: fmt.Sprintf("%v", "otp code cannot be empty")}, nil
	}
	waClient, waClientErr := wa_service.GetWaClient()
	if waClientErr != nil {
		return &wa_grpc_service.SendOtpResponse{StatusCode: "500", Success: false, Data: fmt.Sprintf("%v", waClientErr)}, nil
	}
	err := waClient.SendMessage(otpMessage.PhoneNumber, wa_service.GetOtpMessage(otpMessage.OtpCode))
	if err != nil {
		return &wa_grpc_service.SendOtpResponse{StatusCode: "500", Success: false, Data: fmt.Sprintf("%v", err)}, nil
	}
	return &wa_grpc_service.SendOtpResponse{StatusCode: "200", Success: true, Data: ""}, nil
}
