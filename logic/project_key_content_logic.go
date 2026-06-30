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

var projectKeyContentLogicInstance *ProjectKeyContentLogic
var onceForProjectKeyContentLogicInstance sync.Once = sync.Once{}

type ProjectKeyContentLogic struct {
	dbEngine             *xorm.Engine
	projectKeyContentDAO *dao.ProjectKeyContentDAO
}

func NewProjectKeyContentLogic(ctx context.Context, dbEngine *xorm.Engine,
	projectKeyContentDAO *dao.ProjectKeyContentDAO) *ProjectKeyContentLogic {
	onceForProjectKeyContentLogicInstance.Do(func() {
		projectKeyContentLogicInstance = &ProjectKeyContentLogic{
			dbEngine:             dbEngine,
			projectKeyContentDAO: projectKeyContentDAO,
		}
	})

	return projectKeyContentLogicInstance
}

func (logic *ProjectKeyContentLogic) CreateKeyContent(ctx context.Context, projectEntity *entity.ProjectEntity, key string, content string) (projectKeyContentEntity *entity.ProjectKeyContentEntity, err error) {
	slog.InfoContext(ctx, "create project key-content", slog.Any("projectEntity", projectEntity), slog.String("key", key), slog.String("content", content))
	now := time.Now()
	projectKeyContentEntity = &entity.ProjectKeyContentEntity{
		ProjectId:         projectEntity.Id,
		ProjectKey:        key,
		ProjectKeyContent: content,
		CreateTime:        now,
		UpdateTime:        now,
	}
	err = wrapper.TransactionWrapper(ctx, logic.dbEngine, func(session *xorm.Session) error {
		projectKeyContentEntity, err = logic.projectKeyContentDAO.Insert(ctx, session, projectKeyContentEntity)
		return err
	})
	if err != nil {
		slog.ErrorContext(ctx, "create project key-content", slog.Any("projectEntity", projectEntity), slog.String("key", key), slog.String("content", content), slog.Any("error", err))
		return nil, err
	}

	slog.ErrorContext(ctx, "create project key-content done", slog.Any("projectEntity", projectEntity), slog.String("key", key), slog.String("content", content), slog.Any("projectKeyContentEntity", projectKeyContentEntity))
	return projectKeyContentEntity, nil
}

func (logic *ProjectKeyContentLogic) SaveKeyContent(ctx context.Context, projectEntity *entity.ProjectEntity, key string, content string) (projectKeyContentEntity *entity.ProjectKeyContentEntity, err error) {
	slog.InfoContext(ctx, "save project key-content", slog.Any("projectEntity", projectEntity), slog.String("key", key), slog.String("content", content))

	err = wrapper.TransactionWrapper(ctx, logic.dbEngine, func(session *xorm.Session) error {
		projectKeyContentEntity, err = logic.projectKeyContentDAO.QueryByProjectKey(ctx, session, projectEntity.Id, key)
		if err != nil {
			return err
		}
		if projectKeyContentEntity == nil {
			slog.InfoContext(ctx, "project key-content not found, then create a new entity", slog.Any("projectEntity", projectEntity), slog.String("key", key))

			now := time.Now()
			projectKeyContentEntity = &entity.ProjectKeyContentEntity{
				ProjectId:         projectEntity.Id,
				ProjectKey:        key,
				ProjectKeyContent: content,
				CreateTime:        now,
				UpdateTime:        now,
			}

			projectKeyContentEntity, err = logic.projectKeyContentDAO.Insert(ctx, session, projectKeyContentEntity)
			return err
		}

		projectKeyContentEntity.ProjectKeyContent = content
		projectKeyContentEntity.UpdateTime = time.Now()
		err = logic.projectKeyContentDAO.UpdateByProjectKey(ctx, session, projectKeyContentEntity)
		return err
	})

	if err != nil {
		slog.ErrorContext(ctx, "save project key-content failed", slog.Any("projectEntity", projectEntity), slog.String("key", key), slog.String("content", content), slog.Any("error", err))
		return nil, err
	}

	slog.InfoContext(ctx, "save project key-content done", slog.Any("projectEntity", projectEntity), slog.String("key", key), slog.String("content", content), slog.Any("projectKeyContentEntity", projectKeyContentEntity))
	return projectKeyContentEntity, nil
}

func (logic *ProjectKeyContentLogic) FindByProjectKey(ctx context.Context, projectEntity *entity.ProjectEntity, key string) (projectKeyContentEntity *entity.ProjectKeyContentEntity, err error) {
	slog.InfoContext(ctx, "find project key-content", slog.Any("projectEntity", projectEntity), slog.String("key", key))

	err = wrapper.TransactionWrapper(ctx, logic.dbEngine, func(session *xorm.Session) error {
		projectKeyContentEntity, err = logic.projectKeyContentDAO.QueryByProjectKey(ctx, session, projectEntity.Id, key)
		return err
	})
	if err != nil {
		slog.ErrorContext(ctx, "find project key-content failed", slog.Any("projectEntity", projectEntity), slog.String("key", key), slog.Any("error", err))
		return nil, err
	}
	slog.InfoContext(ctx, "project key-content found", slog.Any("projectEntity", projectEntity), slog.String("key", key))
	return projectKeyContentEntity, nil
}
