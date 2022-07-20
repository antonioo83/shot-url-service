package grpc

import (
	"context"
	"fmt"
	"github.com/antonioo83/shot-url-service/config"
	authInterfaces "github.com/antonioo83/shot-url-service/internal/handlers/auth/interfaces"
	genInterfaces "github.com/antonioo83/shot-url-service/internal/handlers/generators/interfaces"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	pb "github.com/antonioo83/shot-url-service/internal/server/grpc/proto"
	"github.com/antonioo83/shot-url-service/internal/services"
	"google.golang.org/grpc/peer"
	"net/http"
)

// ShortURLServer поддерживает все необходимые методы сервера.
type ShortURLServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedShortURLServer
	Config             config.Config
	ShotURLRepository  interfaces.ShotURLRepository
	UserRepository     interfaces.UserRepository
	DatabaseRepository interfaces.DatabaseRepository
	UserAuthHandler    authInterfaces.UserAuthHandler
	Generator          genInterfaces.ShortLinkGenerator
}

func (s *ShortURLServer) CreateShortURL(ctx context.Context, in *pb.ShortURLRequest) (*pb.ShortURLResponse, error) {
	response := pb.ShortURLResponse{}
	user, err := s.UserRepository.FindByCode(int(in.UserCode))
	if err != nil || user == nil {
		return &response, fmt.Errorf("i can't find user: %w", err)
	}

	p, ok := peer.FromContext(ctx)
	if !ok {
		return &response, fmt.Errorf("i can't get host: %w", err)
	}

	var createShortURLs []services.CreateShortURL
	createShortURLs = append(createShortURLs, services.CreateShortURL{OriginalURL: in.URL, CorrelationID: ""})

	var param services.ShortURLParameters
	param.Config = s.Config
	param.Repository = s.ShotURLRepository
	param.UserRepository = s.UserRepository
	param.Generator = s.Generator
	param.Host = p.Addr.String()
	param.User = *user
	param.CreateShortURLs = &createShortURLs
	result, err := services.SaveShortURLs(param)
	if err != nil {
		return &response, fmt.Errorf("i can't save short url: %w", err)
	}

	response = pb.ShortURLResponse{
		Status:        201,
		CorrelationID: result.ShortURLResponses[0].CorrelationID,
		ShortURL:      result.ShortURLResponses[0].ShortURL,
	}

	return &response, err
}

func (s *ShortURLServer) CreateBatchShortURL(ctx context.Context, in *pb.BatchShortURLRequests) (*pb.BatchShortURLResponses, error) {
	var responses pb.BatchShortURLResponses
	user, err := s.UserRepository.FindByCode(int(in.UserCode))
	if err != nil || user == nil {
		return &responses, fmt.Errorf("i can't find user: %w", err)
	}

	p, ok := peer.FromContext(ctx)
	if !ok {
		return &responses, fmt.Errorf("i can't get host: %w", err)
	}

	var createShortURLs []services.CreateShortURL
	for _, item := range in.Items {
		createShortURLs = append(createShortURLs, services.CreateShortURL{OriginalURL: item.OriginalURL, CorrelationID: item.CorrelationID})
	}

	var param services.ShortURLParameters
	param.Config = s.Config
	param.Repository = s.ShotURLRepository
	param.UserRepository = s.UserRepository
	param.Generator = s.Generator
	param.Host = p.Addr.String()
	param.User = *user
	param.CreateShortURLs = &createShortURLs
	result, err := services.SaveShortURLs(param)
	if err != nil {
		return &responses, fmt.Errorf("i can't save short url: %w", err)
	}

	responses.Status = 201
	for _, res := range result.ShortURLResponses {
		responses.Items = append(responses.Items, &pb.BatchShortURLResponse{ShortURL: res.ShortURL, CorrelationID: res.CorrelationID})
	}

	return &responses, err
}

func (s *ShortURLServer) GetShortURL(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	response := pb.GetResponse{}
	result, err := services.GetShortURL(s.ShotURLRepository, in.Code)
	if err != nil {
		return &response, fmt.Errorf("i can't get short url: %w", err)
	}
	response.Status = int32(result.Status)
	response.OriginalURL = result.OriginalURL

	return &response, nil
}

func (s *ShortURLServer) GetUserURLs(ctx context.Context, in *pb.GetUserURLRequest) (*pb.GetUserURLResponses, error) {
	responses := pb.GetUserURLResponses{}
	user, err := s.UserRepository.FindByCode(int(in.UserCode))
	if err != nil || user == nil {
		return &pb.GetUserURLResponses{}, fmt.Errorf("i can't find user: %w", err)
	}

	result, err := services.GetUserShortUrls(s.ShotURLRepository, user.Code)
	if err != nil {
		return &pb.GetUserURLResponses{}, fmt.Errorf("i can't get user urls: %w", err)
	}

	for _, model := range *result.Models {
		responses.Items = append(responses.Items, &pb.GetUserURLResponse{ShortURL: model.ShortURL, OriginalURL: model.OriginalURL})
	}

	return &responses, nil
}

func (s *ShortURLServer) DeleteUserURLs(ctx context.Context, in *pb.DeleteUserURLRequest) (*pb.DeleteUserURLResponse, error) {
	user, err := s.UserRepository.FindByCode(int(in.UserCode))
	if err != nil || user == nil {
		return nil, fmt.Errorf("i can't find user: %w", err)
	}

	jobCh := make(chan services.ShotURLDelete)
	services.RunDeleteShortURLWorker(jobCh, s.ShotURLRepository, s.Config.DeleteShotURL.WorkersCount)
	services.SendCodesForDeleteToChanel(
		jobCh,
		services.ShotURLDelete{UserCode: user.Code, Codes: in.Codes},
		s.Config.DeleteShotURL.ChunkLength,
	)

	return &pb.DeleteUserURLResponse{Status: int32(http.StatusAccepted)}, nil
}
