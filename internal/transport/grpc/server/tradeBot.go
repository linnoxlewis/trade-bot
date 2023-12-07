package server

import (
	"context"
	pb "github.com/linnoxlewis/trade-bot/internal/transport/grpc/pb/trade-bot"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TradeBotServer struct {
	//orderSrv service.
	pb.UnimplementedPaySystemServiceServer
}

func (t TradeBotServer) mustEmbedUnimplementedPaySystemServiceServer() {
	//TODO implement me
	panic("implement me")
}

func NewTradeBotServer() *TradeBotServer {
	//func NewPaySystemServiceServer(orderSrv service.OrderInterface) *PaySystemServiceServer {
	return &TradeBotServer{
		//	orderSrv: orderSrv,
	}
}

func (t TradeBotServer) Ping(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
