package dao

import (
	"context"
	"sync"
)

var projectKeyDAOInstance *ProjectKeyDAO
var onceForProjecKeytDAOInstance sync.Once = sync.Once{}

type ProjectKeyDAO struct {
}

func NewProjectKeyDAO(ctx context.Context) *ProjectKeyDAO {
	onceForProjecKeytDAOInstance.Do(func() {
		projectKeyDAOInstance = &ProjectKeyDAO{}
	})

	return projectKeyDAOInstance
}
