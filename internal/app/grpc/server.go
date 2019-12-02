package grpc

import (
	"context"
	"net"
	"strconv"

	"github.com/DanielTitkov/antibruteforce-microservice/api"
	"github.com/DanielTitkov/antibruteforce-microservice/internal/app/config"
	"github.com/DanielTitkov/antibruteforce-microservice/internal/bucketstorage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	logger *zap.SugaredLogger
	config *config.AppConfig
	bs     *bucketstorage.BucketStorage
}

func (srv *GRPCServer) Attempt(ctx context.Context, req *api.AttemptRequest) (*api.AttemptResponse, error) {
	srv.logger.Infof("Recieved Attepmt request: %v", req)

	return &api.AttemptResponse{Status: "success", Ok: true}, nil
}

func (srv *GRPCServer) AddToBlacklist(ctx context.Context, req *api.AddToBlacklistRequest) (*api.AddToBlacklistResponse, error) {
	srv.logger.Info("Recieved Add To Blacklist request: %v", req)
	return &api.AddToBlacklistResponse{Status: "success"}, nil
} 

func (srv *GRPCServer) RemoveFromBlacklist(ctx context.Context, req *api.RemoveFromBlacklistRequest) (*api.RemoveFromBlacklistResponse, error) {
	srv.logger.Info("Recieved Remove From Blacklist request: %v", req)
	return &api.RemoveFromBlacklistResponse{Status: "success"}, nil
}

func (srv *GRPCServer) AddToWhitelist(ctx context.Context, req *api.AddToWhitelistRequest) (*api.AddToWhitelistResponse, error) {
	srv.logger.Info("Recieved Add To Whitelist request: %v", req)
	return &api.AddToWhitelistResponse{Status: "success"}, nil
}

func (srv *GRPCServer) RemoveFromWhitelist(ctx context.Context, req *api.RemoveFromWhitelistRequest) (*api.RemoveFromWhitelistResponse, error) {
	srv.logger.Info("Recieved Remove From Whitelist request: %v", req)
	return &api.RemoveFromWhitelistResponse{Status: "success"}, nil
}

func (srv *GRPCServer) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", srv.config.GRPC.Host+":"+strconv.Itoa(srv.config.GRPC.Port))
	if err != nil {
		return err
	}
	srv.logger.Infof("GRPC Server started. Listening on %s:%d", srv.config.GRPC.Host, srv.config.GRPC.Port)
	grpcServer := grpc.NewServer()
	api.RegisterABServiceServer(grpcServer, srv)
	err = grpcServer.Serve(lis)

	return err
}

// New creates new server struct
func New(
	logger *zap.SugaredLogger,
	config *config.AppConfig,
	bs *bucketstorage.BucketStorage,
) *GRPCServer {
	return &GRPCServer{logger, config, bs}
}
