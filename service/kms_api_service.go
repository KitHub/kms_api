package service

import (
	"context"
	"log/slog"
	"sync"

	"github.com/KitHub/kms_api/logic"
	"github.com/KitHub/protocols/kms_api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	kmsapiServiceInstance *KMSAPIService
	kmsapiServiceOnce     sync.Once
)

type KMSAPIService struct {
	kms_api.UnimplementedKMSAPIServer
	projectLogic           *logic.ProjectLogic
	projectKeyContentLogic *logic.ProjectKeyContentLogic
	projectTokenLogic      *logic.ProjectTokenLogic
}

// Load implements [kms_api.KMSAPIServer].
func (s *KMSAPIService) Load(ctx context.Context, req *kms_api.LoadRequest) (rsp *kms_api.LoadResponse, err error) {
	// todo: add cached to improve performance
	slog.InfoContext(ctx, "load key-content", slog.Any("project_id", req.GetProjectId()), slog.Any("key", req.GetKey()))

	projectEntity, err := s.projectLogic.FindById(ctx, req.GetProjectId())
	if err != nil {
		slog.ErrorContext(ctx, "find project failed", slog.Any("projectId", req.GetProjectId()), slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "server error")
	}
	if projectEntity == nil {
		slog.ErrorContext(ctx, "project not found", slog.Any("projectId", req.GetProjectId()))
		return nil, status.Errorf(codes.InvalidArgument, "project not found")
	}

	projectKeyContentEntity, err := s.projectKeyContentLogic.FindByProjectKey(ctx, projectEntity, req.GetKey())
	if err != nil {
		slog.ErrorContext(ctx, "find project key-content failed", slog.Any("projectId", req.GetProjectId()), slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "server error")
	}
	if projectKeyContentEntity == nil {
		rsp = &kms_api.LoadResponse{
			ErrCode: 0,
			ErrMsg:  "ok",
			Data: &kms_api.LoadResponseData{
				Content: "",
			},
		}
		return rsp, nil
	}

	rsp = &kms_api.LoadResponse{
		ErrCode: 0,
		ErrMsg:  "ok",
		Data: &kms_api.LoadResponseData{
			Content: projectKeyContentEntity.ProjectKeyContent,
		},
	}

	slog.InfoContext(ctx, "load key-content done", slog.Any("project_id", req.GetProjectId()), slog.Any("key", req.GetKey()), slog.Any("content", projectKeyContentEntity.ProjectKeyContent))
	return rsp, nil
}

// Store implements [kms_api.KMSAPIServer].
func (s *KMSAPIService) Store(ctx context.Context, req *kms_api.StoreRequest) (rsp *kms_api.StoreResponse, err error) {
	slog.InfoContext(ctx, "store key-content", slog.Any("project_id", req.GetProjectId()), slog.Any("key", req.GetKey()), slog.Any("content", req.GetContent()))

	projectEntity, err := s.projectLogic.FindById(ctx, req.GetProjectId())
	if err != nil {
		slog.ErrorContext(ctx, "find project failed", slog.Any("projectId", req.GetProjectId()), slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "server error")
	}
	if projectEntity == nil {
		slog.ErrorContext(ctx, "project not found", slog.Any("projectId", req.GetProjectId()))
		return nil, status.Errorf(codes.InvalidArgument, "project not found")
	}

	_, err = s.projectKeyContentLogic.SaveKeyContent(ctx, projectEntity, req.GetKey(), req.GetContent())
	if err != nil {
		slog.ErrorContext(ctx, "save project key-content failed", slog.Any("key", req.GetKey()), slog.Any("content", req.GetContent()), slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "server error")
	}

	rsp = &kms_api.StoreResponse{
		ErrCode: 0,
		ErrMsg:  "ok",
		Data:    &kms_api.StoreResponseData{},
	}

	slog.InfoContext(ctx, "store key-content done", slog.Any("project_id", req.GetProjectId()), slog.Any("key", req.GetKey()), slog.Any("content", req.GetContent()))
	return rsp, nil
}

func NewKMSAPIService(ctx context.Context, projectLogic *logic.ProjectLogic, projectKeyContentLogic *logic.ProjectKeyContentLogic, projectTokenLogic *logic.ProjectTokenLogic) *KMSAPIService {
	kmsapiServiceOnce.Do(func() {
		kmsapiServiceInstance = &KMSAPIService{
			projectLogic:           projectLogic,
			projectKeyContentLogic: projectKeyContentLogic,
			projectTokenLogic:      projectTokenLogic,
		}
	})
	return kmsapiServiceInstance
}
