package service

import (
	"context"
	"log/slog"
	"sync"

	"github.com/KitHub/kms_api/logic"
	"github.com/KitHub/protocols/kms_api"
	"google.golang.org/grpc/codes"
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
		rsp = createPBRspWithPBMessageType[kms_api.LoadResponse](ctx, codes.Internal, nil)
		return rsp, nil
	}
	if projectEntity == nil {
		slog.ErrorContext(ctx, "project not found", slog.Any("projectId", req.GetProjectId()))
		rsp = createPBRspWithPBMessageType[kms_api.LoadResponse](ctx, codes.Internal, nil)
		return rsp, nil
	}

	projectKeyContentEntity, err := s.projectKeyContentLogic.FindByProjectKey(ctx, projectEntity, req.GetKey())
	if err != nil {
		slog.ErrorContext(ctx, "find project key-content failed", slog.Any("projectId", req.GetProjectId()), slog.Any("error", err))
		rsp = createPBRspWithPBMessageType[kms_api.LoadResponse](ctx, codes.Internal, nil)
		return rsp, nil
	}
	if projectKeyContentEntity == nil {
		rsp = createPBRspWithPBMessageType[kms_api.LoadResponse](ctx, codes.OK, &kms_api.LoadResponseData{
			Content: "",
		})
		return rsp, nil
	}

	rsp = createPBRspWithPBMessageType[kms_api.LoadResponse](ctx, codes.OK, &kms_api.LoadResponseData{
		Content: projectKeyContentEntity.ProjectKeyContent,
	})

	slog.InfoContext(ctx, "load key-content done", slog.Any("project_id", req.GetProjectId()), slog.Any("key", req.GetKey()), slog.Any("content", projectKeyContentEntity.ProjectKeyContent))
	return rsp, nil
}

// Store implements [kms_api.KMSAPIServer].
func (s *KMSAPIService) Store(ctx context.Context, req *kms_api.StoreRequest) (rsp *kms_api.StoreResponse, err error) {
	slog.InfoContext(ctx, "store key-content", slog.Any("project_id", req.GetProjectId()), slog.Any("key", req.GetKey()), slog.Any("content", req.GetContent()))

	projectEntity, err := s.projectLogic.FindById(ctx, req.GetProjectId())
	if err != nil {
		slog.ErrorContext(ctx, "find project failed", slog.Any("projectId", req.GetProjectId()), slog.Any("error", err))
		rsp = createPBRspWithPBMessageType[kms_api.StoreResponse](ctx, codes.Internal, nil)
		return rsp, nil
	}
	if projectEntity == nil {
		slog.ErrorContext(ctx, "project not found", slog.Any("projectId", req.GetProjectId()))
		rsp = createPBRspWithPBMessageType[kms_api.StoreResponse](ctx, codes.Internal, nil)
		return rsp, nil
	}

	_, err = s.projectKeyContentLogic.SaveKeyContent(ctx, projectEntity, req.GetKey(), req.GetContent())
	if err != nil {
		slog.ErrorContext(ctx, "save project key-content failed", slog.Any("key", req.GetKey()), slog.Any("content", req.GetContent()), slog.Any("error", err))
		rsp = createPBRspWithPBMessageType[kms_api.StoreResponse](ctx, codes.Internal, nil)
		return rsp, nil
	}

	rsp = createPBRspWithPBMessageType[kms_api.StoreResponse](ctx, codes.OK, &kms_api.StoreResponseData{})

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
