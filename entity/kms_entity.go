package entity

import "time"

type ProjectEntity struct {
	Id          int64     `json:"id" xorm:"id,pk,autoincr"`
	ProjectName string    `json:"project_name" xorm:"project_name,varchar"`
	CreateTime  time.Time `json:"create_time" xorm:"create_time"`
	UpdateTime  time.Time `json:"update_time" xorm:"update_time"`
}

func (p *ProjectEntity) TableName() string {
	return "projects"
}

type ProjectTokenEntity struct {
	Id                     int64     `json:"id" xorm:"id,pk,autoincr"`
	ProjectId              int64     `json:"project_id" xorm:"project_id"`
	ProjectToken           string    `json:"project_token" xorm:"project_token"`
	ProjectTokenExpireTime time.Time `json:"project_token_expire_time" xorm:"project_token_expire_time"`
	CreateTime             time.Time `json:"create_time" xorm:"create_time"`
	UpdateTime             time.Time `json:"update_time" xorm:"update_time"`
}

func (p *ProjectTokenEntity) TableName() string {
	return "project_tokens"
}

type ProjectKeyContentEntity struct {
	Id                int64     `json:"id" xorm:"id,pk,autoincr"`
	ProjectId         int64     `json:"project_id" xorm:"project_id"`
	ProjectKey        string    `json:"project_key" xorm:"project_key"`
	ProjectKeyContent string    `json:"project_key_content" xorm:"project_key_content"`
	CreateTime        time.Time `json:"create_time" xorm:"create_time"`
	UpdateTime        time.Time `json:"update_time" xorm:"update_time"`
}

func (p *ProjectKeyContentEntity) TableName() string {
	return "project_key_contents"
}
