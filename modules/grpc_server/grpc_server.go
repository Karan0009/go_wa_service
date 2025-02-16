package grpc_server

import (
	"fmt"
	"net"

	logging "github.com/Karan0009/go_wa_bot/modules/logger"
	wa_grpc_service "github.com/Karan0009/go_wa_bot/modules/wa_grpc_service"
	wa_grpc_service_proto "github.com/Karan0009/go_wa_bot/modules/wa_grpc_service/proto"
	"google.golang.org/grpc"
)

// InitGrpcServer initializes and starts the gRPC server.
func InitGrpcServer(port string) (*grpc.Server, error) {
	// Define the gRPC server port
	address := fmt.Sprintf(":%s", port)

	// Start listening on the specified port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %s: %v", port, err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()
	wa_grpc_service_proto.RegisterWAServiceServer(grpcServer, wa_grpc_service.GetMyWaGrpcServiceServer())
	// Register services here (e.g., pb.RegisterYourServiceServer(grpcServer, &YourService{}))

	logging.NewLogger("grpc_server").Info(fmt.Sprintf("gRPC server started on port %s", port))

	// Start serving requests
	if err := grpcServer.Serve(listener); err != nil {
		return nil, fmt.Errorf("failed to start gRPC server: %v", err)
	}

	return grpcServer, nil
}
