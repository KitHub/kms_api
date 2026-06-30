package logic

import (
	"context"
	"sync"

	"github.com/KitHub/kms_api/dao"
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
