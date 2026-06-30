package dao

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/KitHub/kms_api/entity"
	"xorm.io/xorm"
)

var projectDAOInstance *ProjectDAO
var onceForProjectDAOInstance sync.Once = sync.Once{}

type ProjectDAO struct {
}

func NewProjectDAO(ctx context.Context) *ProjectDAO {
	onceForProjectDAOInstance.Do(func() {
		projectDAOInstance = &ProjectDAO{}
	})

	return projectDAOInstance
}

func (dao *ProjectDAO) Insert(ctx context.Context, session *xorm.Session, projectEntity *entity.ProjectEntity) (*entity.ProjectEntity, error) {
	slog.InfoContext(ctx, "insert project", slog.Any("projectEntity", projectEntity))
	rowsEffected, err := session.Insert(projectEntity)
	if err != nil {
		slog.ErrorContext(ctx, "insert project failed",
			slog.Any("project", projectEntity), slog.Any("error", err))
		return nil, err
	}
	if rowsEffected == 0 {
		errMsg := "no rows affected when inserting project"
		err = errors.New(errMsg)
		slog.ErrorContext(ctx, errMsg,
			slog.Any("project", projectEntity.Id))
		return nil, err
	}
	slog.InfoContext(ctx, "project inserted", slog.Any("project", projectEntity),
		slog.Any("rows_affected", rowsEffected))

	return projectEntity, nil
}

func (dao *ProjectDAO) QueryById(ctx context.Context,
	session *xorm.Session, id int64) (*entity.ProjectEntity, error) {
	projectEntity := &entity.ProjectEntity{}
	has, err := session.Where("id = ?", id).Get(projectEntity)
	if err != nil {
		slog.ErrorContext(ctx, "query project by id failed", slog.Int64("project_id", id), slog.Any("error", err))
		return nil, err
	}
	if !has {
		slog.InfoContext(ctx, "project not found", slog.Int64("project_id", id))
		return nil, nil
	}
	slog.InfoContext(ctx, "project found", slog.Any("project", projectEntity))
	return projectEntity, nil
}

func (dao *ProjectDAO) QueryByProjectName(ctx context.Context, session *xorm.Session, projectName string) (*entity.ProjectEntity, error) {
	projectEntity := &entity.ProjectEntity{}
	has, err := session.Where("project_name = ?", projectName).Get(projectEntity)
	if err != nil {
		slog.ErrorContext(ctx, "query project by project_name failed", slog.String("project_name", projectName), slog.Any("error", err))
		return nil, err
	}
	if !has {
		slog.InfoContext(ctx, "project not found", slog.String("project_name", projectName))
		return nil, nil
	}
	slog.InfoContext(ctx, "project found", slog.Any("project", projectEntity))
	return projectEntity, nil
}
