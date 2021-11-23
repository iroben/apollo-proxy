package model

import (
	"crypto/md5"
	"fmt"
	"log"
	"time"
)

type Project struct {
	BaseModel
	Id          int        `gorm:"column:id;primaryKey"`
	Hash        string     `gorm:"column:hash"`
	Env         string     `gorm:"column:env" uri:"env"`
	Project     string     `gorm:"column:project" uri:"project"`
	Branch      string     `gorm:"column:branch" uri:"branch"`
	ApolloAppId string     `gorm:"column:apollo_app_id" uri:"apollo_app_id"`
	ClusterName string     `gorm:"column:cluster_name" uri:"cluster_name"`
	Namespace   string     `gorm:"column:namespace" uri:"namespace"`
	Created     TimeYmdHis `gorm:"column:created"`
}

func (Project) TableName() string {
	return "project"
}
func (m *Project) Save() bool {
	/**
	开发环境的配置信息不保存
	norecord 项目也不保存
	*/
	if m.Env == "dev" || m.Project == "norecord" {
		return true
	}
	hashStr := fmt.Sprintf("project=%s AND branch=%s AND apollo_app_id=%s AND cluster_name=%s AND namespace=%s AND env=%s",
		m.Project,
		m.Branch,
		m.ApolloAppId,
		m.ClusterName,
		m.Namespace,
		m.Env)
	h := md5.New()
	h.Write([]byte(hashStr))
	m.Hash = fmt.Sprintf("%X", h.Sum(nil))
	m.Created = TimeYmdHis(time.Now())
	err := DB.Where("hash=?", m.Hash).FirstOrCreate(m).Error
	if err != nil {
		log.Println("保存失败：", err)
		return false
	}
	return true

}

/**
获取所有的jenkins项目
*/
func JenkinsApps() []*Project {
	var retVal []*Project
	err := DB.Model(&Project{}).Group("project,env").Find(&retVal).Error
	if err != nil {
		log.Println("JenkinsApps Error: ", err)
		return nil
	}
	return retVal
}

func (m *Project) FindProjects() []*Project {
	var retVal []*Project
	err := DB.Where("project=? AND env=?", m.Project, m.Env).Find(&retVal).Error
	if err != nil {
		log.Println("FindProjectByApolloId Error: ", err)
		return nil
	}
	return retVal
}
