package logic

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/KitHub/kms_api/dao"
	"github.com/KitHub/kms_api/entity"
	"github.com/KitHub/kms_api/wrapper"
	"xorm.io/xorm"
)

var projectLogicInstance *ProjectLogic
var onceForProjectLogicInstance sync.Once = sync.Once{}

type ProjectLogic struct {
	dbEngine   *xorm.Engine
	projectDAO *dao.ProjectDAO
}

func NewProjectLogic(ctx context.Context, dbEngine *xorm.Engine,
	projectDAO *dao.ProjectDAO) *ProjectLogic {
	onceForProjectLogicInstance.Do(func() {
		projectLogicInstance = &ProjectLogic{
			dbEngine:   dbEngine,
			projectDAO: projectDAO,
		}
	})

	return projectLogicInstance
}

func (logic *ProjectLogic) CreateProject(ctx context.Context, projectName string) (projectEntity *entity.ProjectEntity, err error) {
	slog.InfoContext(ctx, "create project", slog.String("projectName", projectName))
	now := time.Now()
	projectEntity = &entity.ProjectEntity{
		ProjectName: projectName,
		CreateTime:  now,
		UpdateTime:  now,
	}
	err = wrapper.TransactionWrapper(ctx, logic.dbEngine, func(session *xorm.Session) error {
		projectEntity, err = logic.projectDAO.Insert(ctx, session, projectEntity)
		return err
	})
	if err != nil {
		slog.ErrorContext(ctx, "execute transaction failed when creating project", slog.Any("projectEntity", projectEntity), slog.Any("error", err))
		return nil, err
	}
	slog.InfoContext(ctx, "create project done", slog.String("projectName", projectName), slog.Any("projectEntity", projectEntity))
	return projectEntity, nil
}

func (logic *ProjectLogic) FindByName(ctx context.Context, projectName string) (projectEntity *entity.ProjectEntity, err error) {
	slog.InfoContext(ctx, "find project", slog.String("projectName", projectName))
	projectEntity = &entity.ProjectEntity{
		ProjectName: projectName,
	}
	err = wrapper.TransactionWrapper(ctx, logic.dbEngine, func(session *xorm.Session) error {
		projectEntity, err = logic.projectDAO.QueryByProjectName(ctx, session, projectName)
		return err
	})
	if err != nil {
		slog.ErrorContext(ctx, "find project failed", slog.String("projectName", projectName), slog.Any("error", err))
		return nil, err
	}
	slog.InfoContext(ctx, "project found", slog.String("projectName", projectName), slog.Any("projectEntity", projectEntity))
	return projectEntity, nil
}
