package logic

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/KitHub/kms_api/config"
	"github.com/KitHub/kms_api/dao"
	"github.com/KitHub/kms_api/entity"
	"github.com/KitHub/kms_api/wrapper"
	"github.com/google/uuid"
	"xorm.io/xorm"
)

var projectTokenLogicInstance *ProjectTokenLogic
var onceForProjectTokenLogicInstance sync.Once = sync.Once{}

type ProjectTokenLogic struct {
	projectTokensConfig *config.ProjectTokensConfigEntity
	dbEngine            *xorm.Engine
	projectTokenDAO     *dao.ProjectTokenDAO
}

func NewProjectTokenLogic(ctx context.Context, projectTokensConfig *config.ProjectTokensConfigEntity, dbEngine *xorm.Engine,
	projectTokenDAO *dao.ProjectTokenDAO) *ProjectTokenLogic {
	onceForProjectTokenLogicInstance.Do(func() {
		projectTokenLogicInstance = &ProjectTokenLogic{
			projectTokensConfig: projectTokensConfig,
			dbEngine:            dbEngine,
			projectTokenDAO:     projectTokenDAO,
		}
	})

	return projectTokenLogicInstance
}

func (logic *ProjectTokenLogic) Create(ctx context.Context, projectEntity *entity.ProjectEntity) (projectTokenEntity *entity.ProjectTokenEntity, err error) {
	slog.InfoContext(ctx, "create project token", slog.Any("projectEntity", projectEntity))
	now := time.Now()
	token := generateToken(ctx)
	projectTokenEntity = &entity.ProjectTokenEntity{
		ProjectId:              projectEntity.Id,
		ProjectToken:           token,
		ProjectTokenExpireTime: now.Add(time.Duration(logic.projectTokensConfig.DefaultDurationInDaysForNewToken) * time.Hour * 24),
		CreateTime:             now,
		UpdateTime:             now,
	}

	err = wrapper.TransactionWrapper(ctx, logic.dbEngine, func(session *xorm.Session) error {
		projectTokenEntity, err = logic.projectTokenDAO.Insert(ctx, session, projectTokenEntity)
		return err
	})
	if err != nil {
		slog.ErrorContext(ctx, "create project token failed", slog.Any("projectEntity", projectEntity), slog.String("token", token), slog.Any("error", err))
		return nil, err
	}

	slog.InfoContext(ctx, "create project token done", slog.Any("projectEntity", projectEntity), slog.String("token", token), slog.Any("projectTokenEntity", projectTokenEntity))
	return projectTokenEntity, nil
}

func (logic *ProjectTokenLogic) Find(ctx context.Context, projectEntity *entity.ProjectEntity, token string) (projectTokenEntity *entity.ProjectTokenEntity, err error) {
	slog.InfoContext(ctx, "find project token", slog.Any("projectEntity", projectEntity), slog.String("token", token))

	err = wrapper.TransactionWrapper(ctx, logic.dbEngine, func(session *xorm.Session) error {
		projectTokenEntity, err = logic.projectTokenDAO.QueryByProjectToken(ctx, session, projectEntity.Id, token)
		return err
	})
	if err != nil {
		slog.ErrorContext(ctx, "find project token failed", slog.Any("projectEntity", projectEntity), slog.String("token", token), slog.Any("error", err))
		return nil, err
	}

	slog.InfoContext(ctx, "find project token done", slog.Any("projectEntity", projectEntity), slog.String("token", token), slog.Any("projectTokenEntity", projectTokenEntity))
	return projectTokenEntity, nil
}

func (logic *ProjectTokenLogic) Disable(ctx context.Context, projectEntity *entity.ProjectEntity, token string) (err error) {
	slog.InfoContext(ctx, "disable project token", slog.Any("projectEntity", projectEntity), slog.String("token", token))

	var projectTokenEntity *entity.ProjectTokenEntity
	err = wrapper.TransactionWrapper(ctx, logic.dbEngine, func(session *xorm.Session) error {
		projectTokenEntity, err = logic.projectTokenDAO.QueryByProjectToken(ctx, session, projectEntity.Id, token)
		if err != nil {
			return err
		}
		if projectTokenEntity == nil {
			errMsg := "project token not found"
			err = errors.New(errMsg)
			slog.ErrorContext(ctx, errMsg, slog.Any("projectEntity", projectEntity), slog.String("token", token))
		}

		projectTokenEntity.ProjectTokenExpireTime = time.Now().Add(time.Duration(logic.projectTokensConfig.DefaultDurationInDaysForUselessToken) * time.Hour * 24)

		err = logic.projectTokenDAO.UpdateByProjectToken(ctx, session, projectTokenEntity)
		return err
	})
	if err != nil {
		slog.ErrorContext(ctx, "disable project token failed", slog.Any("projectEntity", projectEntity), slog.String("token", token), slog.Any("error", err))
		return err
	}

	slog.InfoContext(ctx, "disable project token done", slog.Any("projectEntity", projectEntity), slog.String("token", token), slog.Any("projectTokenEntity", projectTokenEntity))
	return nil
}

func (logic *ProjectTokenLogic) Renew(ctx context.Context, projectEntity *entity.ProjectEntity, token string) (newProjectTokenEntity *entity.ProjectTokenEntity, err error) {
	slog.InfoContext(ctx, "renew project token", slog.Any("projectEntity", projectEntity), slog.String("token", token))

	var projectTokenEntity *entity.ProjectTokenEntity
	err = wrapper.TransactionWrapper(ctx, logic.dbEngine, func(session *xorm.Session) error {
		projectTokenEntity, err = logic.projectTokenDAO.QueryByProjectToken(ctx, session, projectEntity.Id, token)
		if err != nil {
			return err
		}
		if projectTokenEntity == nil {
			errMsg := "project token not found"
			err = errors.New(errMsg)
			slog.ErrorContext(ctx, errMsg, slog.Any("projectEntity", projectEntity), slog.String("token", token))
		}

		now := time.Now()

		projectTokenEntity.ProjectTokenExpireTime = now.Add(time.Duration(logic.projectTokensConfig.DefaultDurationInDaysForUselessToken) * time.Hour * 24)
		err = logic.projectTokenDAO.UpdateByProjectToken(ctx, session, projectTokenEntity)

		newProjectTokenEntity = &entity.ProjectTokenEntity{
			ProjectId:              projectEntity.Id,
			ProjectToken:           generateToken(ctx),
			ProjectTokenExpireTime: now.Add(time.Duration(logic.projectTokensConfig.DefaultDurationInDaysForNewToken) * time.Hour * 24),
			CreateTime:             now,
			UpdateTime:             now,
		}
		newProjectTokenEntity, err = logic.projectTokenDAO.Insert(ctx, session, newProjectTokenEntity)
		return err
	})
	if err != nil {
		slog.ErrorContext(ctx, "renew project token failed", slog.Any("projectEntity", projectEntity), slog.String("token", token), slog.Any("error", err))
		return nil, err
	}

	slog.InfoContext(ctx, "renew project token done", slog.Any("projectEntity", projectEntity), slog.String("token", token), slog.Any("projectTokenEntity", projectTokenEntity))
	return newProjectTokenEntity, nil
}

func generateToken(ctx context.Context) string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
