package grpcClient

import (
	"fmt"
	"post-service/config"
	pbp "post-service/genproto/customer"
	pbr "post-service/genproto/reyting"

	"google.golang.org/grpc"
)

// GrpcClientI ...
type GrpcClientI interface {
	Customer() pbp.CustomerServiceClient
	Ranking() pbr.RankingServiceClient
}

// GrpcClient ...
type GrpcClient struct {
	cfg             config.Config
	customerService pbp.CustomerServiceClient
	rankingService  pbr.RankingServiceClient
}

// New ...
func New(cfg config.Config) (*GrpcClient, error) {

	con, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.CustomerServiceHost, cfg.CustomerServicePort),
		grpc.WithInsecure())

	if err != nil {
		return nil, fmt.Errorf("customer service dial host:%s port:%d", cfg.CustomerServiceHost, cfg.CustomerServicePort)
	}

	conRanking, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.RankingServiceHost, cfg.RankingServicePort),
		grpc.WithInsecure())

	return &GrpcClient{
		cfg:             cfg,
		customerService: pbp.NewCustomerServiceClient(con),
		rankingService:  pbr.NewRankingServiceClient(conRanking),
	}, nil
}

func (g *GrpcClient) Customer() pbp.CustomerServiceClient {
	return g.customerService
}

func (g *GrpcClient) Ranking() pbr.RankingServiceClient {
	return g.rankingService
}
