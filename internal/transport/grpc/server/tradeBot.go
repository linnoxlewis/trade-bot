package server

import (
	"context"
	pb "github.com/linnoxlewis/trade-bot/internal/transport/grpc/pb/trade-bot"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TradeBotServer struct {
	pb.UnimplementedPaySystemServiceServer
}

func NewTradeBotServer() *TradeBotServer {
	return &TradeBotServer{}
}

func (t TradeBotServer) Ping(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
