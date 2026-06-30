package servicecontext

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/KitHub/kms_api/component"
	"github.com/KitHub/kms_api/config"
	"github.com/KitHub/kms_api/dao"
	"github.com/KitHub/kms_api/logic"
	"github.com/KitHub/kms_api/service"
	"gopkg.in/natefinch/lumberjack.v2"
	"xorm.io/xorm"
)

type ServiceContext struct {
	Logger                 *slog.Logger
	CronComponent          *component.CronComponent
	InitComponent          *component.InitComponent
	ShutdownComponent      *component.ShutdownComponent
	DBEngine               *xorm.Engine
	ProjectDAO             *dao.ProjectDAO
	ProjectKeyContentDAO   *dao.ProjectKeyContentDAO
	ProjectTokenDAO        *dao.ProjectTokenDAO
	ProjectLogic           *logic.ProjectLogic
	ProjectKeyContentLogic *logic.ProjectKeyContentLogic
	ProjectTokenLogic      *logic.ProjectTokenLogic
	KMSAPIService          *service.KMSAPIService
}

var gServiceCtx *ServiceContext
var once sync.Once

func InitServiceContext(ctx context.Context, configEntity *config.ConfigEntity) (
	serviceCtx *ServiceContext, err error) {
	slog.InfoContext(ctx, "init service context")

	once.Do(func() {
		logger, innerErr := initLog(ctx, configEntity.LogConfig)
		if innerErr != nil {
			slog.ErrorContext(ctx, "init log failed", slog.Any("error", innerErr))
			err = innerErr
			return
		}

		cronComponent := component.NewCronConponent()
		initComponent := component.NewInitComponent(ctx)
		shutdownComponent := component.NewShutdownComponent(ctx)

		dbEngine, innerErr := initDB(ctx, configEntity.DBConfig, shutdownComponent)
		if innerErr != nil {
			slog.ErrorContext(ctx, "init database failed", slog.Any("error", innerErr))
			err = innerErr
			return
		}

		projectDAO := dao.NewProjectDAO(ctx)
		projectKeyContentDAO := dao.NewProjectKeyContentDAO(ctx)
		projectTokenDAO := dao.NewProjectTokenDAO(ctx)

		projectLogic := logic.NewProjectLogic(ctx, dbEngine, projectDAO)
		projectKeyContentLogic := logic.NewProjectKeyContentLogic(ctx, dbEngine, projectKeyContentDAO)
		projectTokenLogic := logic.NewProjectTokenLogic(ctx, dbEngine, projectTokenDAO)

		kmsapiService := service.NewKMSAPIService(ctx, projectLogic, projectKeyContentLogic, projectTokenLogic)

		gServiceCtx = &ServiceContext{
			Logger:                 logger,
			ShutdownComponent:      shutdownComponent,
			InitComponent:          initComponent,
			CronComponent:          cronComponent,
			KMSAPIService:          kmsapiService,
			DBEngine:               dbEngine,
			ProjectDAO:             projectDAO,
			ProjectKeyContentDAO:   projectKeyContentDAO,
			ProjectTokenDAO:        projectTokenDAO,
			ProjectLogic:           projectLogic,
			ProjectKeyContentLogic: projectKeyContentLogic,
			ProjectTokenLogic:      projectTokenLogic,
		}
	})

	slog.InfoContext(ctx, "init service context done")
	return gServiceCtx, err
}

func initLog(ctx context.Context, logConfig *config.LogConfigEntity) (
	*slog.Logger, error) {
	log := &lumberjack.Logger{
		Filename:   logConfig.Filename,   // 日志文件路径
		MaxSize:    logConfig.MaxSize,    // 每个日志文件的最大大小（以MB为单位）
		MaxBackups: logConfig.MaxBackups, // 保留旧文件的最大数量
		MaxAge:     logConfig.MaxAge,     // 保留旧文件的最大天数
		Compress:   logConfig.Compress,   // 是否压缩旧文件
		LocalTime:  logConfig.LocalTime,  // 是否使用本地时间戳
	}
	serviceLogger := slog.New(slog.NewTextHandler(log, nil))
	slog.SetDefault(serviceLogger)
	return serviceLogger, nil
}

func initDB(ctx context.Context, dbConfig *config.DBConfigEntity,
	shutdownComponent *component.ShutdownComponent) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine(dbConfig.DriverName, dbConfig.DataSourceName)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to initialize database connection",
			slog.Any("error", err))
		return nil, err
	}

	engine.SetMaxIdleConns(dbConfig.MaxIdleConns)
	engine.SetMaxOpenConns(dbConfig.MaxOpenConns)
	engine.SetConnMaxLifetime(
		time.Duration(dbConfig.ConnMaxLifetime) * time.Second)

	err = engine.PingContext(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to ping database", slog.Any("error", err))
		return nil, err
	}

	shutdownComponent.RegisterShutdownCallback(func(ctx context.Context) error {
		return engine.Close()
	})

	slog.InfoContext(ctx, "Database connection initialized successfully")
	return engine, nil
}

func GetServiceContext() *ServiceContext {
	return gServiceCtx
}
