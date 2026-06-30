package logic

import (
	"context"
	"sync"

	"github.com/KitHub/kms_api/dao"
	"xorm.io/xorm"
)

var projectKeyLogicInstance *ProjectKeyLogic
var onceForProjectKeyLogicInstance sync.Once = sync.Once{}

type ProjectKeyLogic struct {
	dbEngine      *xorm.Engine
	projectKeyDAO *dao.ProjectKeyDAO
}

func NewProjectKeyLogic(ctx context.Context, dbEngine *xorm.Engine,
	projectKeyDAO *dao.ProjectKeyDAO) *ProjectKeyLogic {
	onceForProjectKeyLogicInstance.Do(func() {
		projectKeyLogicInstance = &ProjectKeyLogic{
			dbEngine:      dbEngine,
			projectKeyDAO: projectKeyDAO,
		}
	})

	return projectKeyLogicInstance
}
