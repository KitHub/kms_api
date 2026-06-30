package logic

import (
	"context"
	"sync"

	"github.com/KitHub/kms_api/dao"
	"xorm.io/xorm"
)

var projectTokenLogicInstance *ProjectTokenLogic
var onceForProjectTokenLogicInstance sync.Once = sync.Once{}

type ProjectTokenLogic struct {
	dbEngine        *xorm.Engine
	projectTokenDAO *dao.ProjectTokenDAO
}

func NewProjectTokenLogic(ctx context.Context, dbEngine *xorm.Engine,
	projectTokenDAO *dao.ProjectTokenDAO) *ProjectTokenLogic {
	onceForProjectTokenLogicInstance.Do(func() {
		projectTokenLogicInstance = &ProjectTokenLogic{
			dbEngine:        dbEngine,
			projectTokenDAO: projectTokenDAO,
		}
	})

	return projectTokenLogicInstance
}
