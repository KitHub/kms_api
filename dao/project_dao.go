package dao

import (
	"context"
	"sync"
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
