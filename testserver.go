package emulator

import (
	"net"

	tasks "google.golang.org/genproto/googleapis/cloud/tasks/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TestServer struct {
	listener net.Listener
	server   *grpc.Server
	conn     *grpc.ClientConn
}

func NewTestServer() *TestServer {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	emulatorServer := NewServer()

	tasks.RegisterCloudTasksServer(grpcServer, emulatorServer)

	s := TestServer{
		listener: lis,
		server:   grpcServer,
	}

	go grpcServer.Serve(lis)

	conn, err := grpc.Dial(
		lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	s.conn = conn

	return &s
}

func (s *TestServer) Close() {
	s.server.GracefulStop()
	_ = s.conn.Close()
	_ = s.listener.Close()
}

func (s *TestServer) Address() string {
	return s.listener.Addr().String()
}

func (s *TestServer) Connection() *grpc.ClientConn {
	return s.conn
}
