package dao

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/KitHub/kms_api/entity"
	"xorm.io/xorm"
)

var projectTokenDAOInstance *ProjectTokenDAO
var onceForProjectTokenDAOInstance sync.Once = sync.Once{}

type ProjectTokenDAO struct {
}

func NewProjectTokenDAO(ctx context.Context) *ProjectTokenDAO {
	onceForProjectTokenDAOInstance.Do(func() {
		projectTokenDAOInstance = &ProjectTokenDAO{}
	})

	return projectTokenDAOInstance
}

func (dao *ProjectTokenDAO) Insert(ctx context.Context, session *xorm.Session, projectTokenEntity *entity.ProjectTokenEntity) (*entity.ProjectTokenEntity, error) {
	slog.InfoContext(ctx, "insert project token", slog.Any("projectTokenEntity", projectTokenEntity))
	rowsEffected, err := session.Insert(projectTokenEntity)
	if err != nil {
		slog.ErrorContext(ctx, "insert project token failed",
			slog.Any("projectTokenEntity", projectTokenEntity), slog.Any("error", err))
		return nil, err
	}
	if rowsEffected == 0 {
		errMsg := "no rows affected when inserting project token"
		err = errors.New(errMsg)
		slog.ErrorContext(ctx, errMsg, slog.Any("projectTokenEntity", projectTokenEntity))
		return nil, err
	}
	slog.InfoContext(ctx, "project key-content inserted", slog.Any("projectTokenEntity", projectTokenEntity),
		slog.Any("rows_affected", rowsEffected))

	return projectTokenEntity, nil
}

func (dao *ProjectTokenDAO) UpdateByProjectToken(ctx context.Context, session *xorm.Session, projectTokenEntity *entity.ProjectTokenEntity) error {
	slog.InfoContext(ctx, "update project key-content", slog.Any("projectTokenEntity", projectTokenEntity))
	rowsEffected, err := session.Where("project_id", projectTokenEntity.ProjectId).And("project_token", projectTokenEntity.ProjectToken).Update(projectTokenEntity)
	if err != nil {
		slog.ErrorContext(ctx, "update project token failed",
			slog.Any("projectTokenEntity", projectTokenEntity), slog.Any("error", err))
		return err
	}
	if rowsEffected != 1 {
		errMsg := "no rows affected when update package"
		err = errors.New(errMsg)
		slog.ErrorContext(ctx, errMsg, slog.Any("projectTokenEntity", projectTokenEntity))
		return err
	}
	slog.InfoContext(ctx, "project token updated", slog.Any("projectTokenEntity", projectTokenEntity),
		slog.Any("rows_affected", rowsEffected))

	return nil
}

func (dao *ProjectTokenDAO) QueryByProjectToken(ctx context.Context,
	session *xorm.Session, projectId int64, projectToken string) (*entity.ProjectTokenEntity, error) {
	projectTokenEntity := &entity.ProjectTokenEntity{}
	has, err := session.Where("project_id = ?", projectId).And("project_token", projectToken).Get(projectTokenEntity)
	if err != nil {
		slog.ErrorContext(ctx, "query project token failed", slog.Int64("project_id", projectId), slog.String("project_token", projectToken), slog.Any("error", err))
		return nil, err
	}
	if !has {
		slog.InfoContext(ctx, "project token not found", slog.Int64("project_id", projectId), slog.String("project_token", projectToken))
		return nil, nil
	}

	slog.InfoContext(ctx, "project token found", slog.Int64("project_id", projectId), slog.String("project_token", projectToken), slog.Any("projectTokenEntity", projectTokenEntity))
	return projectTokenEntity, nil
}
