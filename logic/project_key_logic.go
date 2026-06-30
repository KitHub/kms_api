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
	dbEngine             *xorm.Engine
	projectKeyContentDAO *dao.ProjectKeyContentDAO
}

func NewProjectKeyLogic(ctx context.Context, dbEngine *xorm.Engine,
	projectKeyContentDAO *dao.ProjectKeyContentDAO) *ProjectKeyLogic {
	onceForProjectKeyLogicInstance.Do(func() {
		projectKeyLogicInstance = &ProjectKeyLogic{
			dbEngine:             dbEngine,
			projectKeyContentDAO: projectKeyContentDAO,
		}
	})

	return projectKeyLogicInstance
}
