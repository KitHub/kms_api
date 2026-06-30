package dao

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/KitHub/kms_api/entity"
	"xorm.io/xorm"
)

var projectKeyContentDAOInstance *ProjectKeyContentDAO
var onceForProjectKeyContentDAOInstance sync.Once = sync.Once{}

type ProjectKeyContentDAO struct {
}

func NewProjectKeyContentDAO(ctx context.Context) *ProjectKeyContentDAO {
	onceForProjectKeyContentDAOInstance.Do(func() {
		projectKeyContentDAOInstance = &ProjectKeyContentDAO{}
	})

	return projectKeyContentDAOInstance
}

func (dao *ProjectKeyContentDAO) Insert(ctx context.Context, session *xorm.Session, projectKeyContentEntity *entity.ProjectKeyContentEntity) (*entity.ProjectKeyContentEntity, error) {
	slog.InfoContext(ctx, "insert project key-content", slog.Any("projectKeyContentEntity", projectKeyContentEntity))
	rowsEffected, err := session.Insert(projectKeyContentEntity)
	if err != nil {
		slog.ErrorContext(ctx, "insert project key-content failed",
			slog.Any("projectKeyContentEntity", projectKeyContentEntity), slog.Any("error", err))
		return nil, err
	}
	if rowsEffected == 0 {
		errMsg := "no rows affected when inserting project key-content"
		err = errors.New(errMsg)
		slog.ErrorContext(ctx, errMsg, slog.Any("projectKeyContentEntity", projectKeyContentEntity))
		return nil, err
	}
	slog.InfoContext(ctx, "project key-content inserted", slog.Any("projectKeyContentEntity", projectKeyContentEntity),
		slog.Any("rows_affected", rowsEffected))

	return projectKeyContentEntity, nil
}

func (dao *ProjectKeyContentDAO) UpdateByProjectKey(ctx context.Context, session *xorm.Session, projectKeyContentEntity *entity.ProjectKeyContentEntity) error {
	slog.InfoContext(ctx, "update project key-content", slog.Any("projectKeyContentEntity", projectKeyContentEntity))
	rowsEffected, err := session.Where("project_key", projectKeyContentEntity.ProjectKey).Update(projectKeyContentEntity)
	if err != nil {
		slog.ErrorContext(ctx, "update project key-content failed",
			slog.Any("projectKeyContentEntity", projectKeyContentEntity), slog.Any("error", err))
		return err
	}
	if rowsEffected != 1 {
		errMsg := "no rows affected when update project key-content"
		err = errors.New(errMsg)
		slog.ErrorContext(ctx, errMsg, slog.Any("projectKeyContentEntity", projectKeyContentEntity))
		return err
	}
	slog.InfoContext(ctx, "project key-content updated", slog.Any("projectKeyContentEntity", projectKeyContentEntity),
		slog.Any("rows_affected", rowsEffected))

	return nil
}

func (dao *ProjectKeyContentDAO) QueryByProjectKey(ctx context.Context,
	session *xorm.Session, projectId int64, projectKey string) (*entity.ProjectKeyContentEntity, error) {
	projectKeyContentEntity := &entity.ProjectKeyContentEntity{}
	has, err := session.Where("project_id = ?", projectId).And("project_key", projectKey).Get(projectKeyContentEntity)
	if err != nil {
		slog.ErrorContext(ctx, "query project key-content failed", slog.Int64("project_id", projectId), slog.String("project_key", projectKey), slog.Any("error", err))
		return nil, err
	}
	if !has {
		slog.InfoContext(ctx, "project key-content not found", slog.Int64("project_id", projectId), slog.String("project_key", projectKey))
		return nil, nil
	}

	slog.InfoContext(ctx, "project key-content found", slog.Int64("project_id", projectId), slog.String("project_key", projectKey), slog.Any("projectKeyContentEntity", projectKeyContentEntity))
	return projectKeyContentEntity, nil
}
