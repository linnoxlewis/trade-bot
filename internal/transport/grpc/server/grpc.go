package server

import (
	pb "github.com/linnoxlewis/trade-bot/internal/transport/grpc/pb/trade-bot"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	"google.golang.org/grpc"
	"net"
)

type TransportInterface interface {
	StartServer()
	StopServer()
}

type Grpc struct {
	server *grpc.Server
	port   string
	logger *log.Logger
}

func NewGrpc(port string, logger log.Logger) *Grpc {
	srv := grpc.NewServer()
	server := NewTradeBotServer()
	pb.RegisterPaySystemServiceServer(srv, server)

	return &Grpc{server: srv, port: port, logger: &logger}
}

func (g *Grpc) StartServer() {
	g.logger.InfoLog.Println("GRPC Server transport starting...")

	connection, err := net.Listen("tcp", g.port)
	if err != nil {
		g.logger.ErrorLog.Panic(err)
	}

	err = g.server.Serve(connection)
	if err != nil {
		g.logger.ErrorLog.Panic(err)
	}
}

func (g *Grpc) StopServer() {
	g.logger.InfoLog.Println("GRPC Server transport stopping...")
	g.server.Stop()
}
