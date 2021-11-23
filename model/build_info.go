package model

import (
	"log"
	"time"
)

type BuildInfo struct {
	BaseModel
	Id        int        `gorm:"column:id;primaryKey" `
	JobId     int        `gorm:"column:job_id;" `
	Count     int        `gorm:"column:count;" `
	Project   string     `gorm:"column:project;" `
	State     string     `gorm:"column:state;" `
	Namespace string     `gorm:"column:namespace;" `
	Result    string     `gorm:"column:result;" `
	Env       string     `gorm:"column:env;" `
	Branch    string     `gorm:"column:branch;" `
	Created   TimeYmdHis `gorm:"column:created;"`
}

func (BuildInfo) TableName() string {
	return "build_info"
}

func (m *BuildInfo) Save() bool {
	if m.JobId == 0 {
		return false
	}
	m.Created = TimeYmdHis(time.Now())
	err := DB.Where("project=? and job_id=?",
		m.Project, m.JobId).FirstOrCreate(m).Error
	if err != nil {
		log.Println("BuildInfo Save Error: ", err)
		return false
	}
	return true
}

func (m *BuildInfo) Update() bool {
	err := DB.Model(&BuildInfo{}).Where("id=?", m.Id).Updates(m).Error
	if err != nil {
		log.Println("BuildInfo Update Error: ", err)
		return false
	}
	return true
}
