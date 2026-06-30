package dao

import (
	"context"
	"sync"
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
