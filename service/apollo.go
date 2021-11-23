package service

import (
	"apollo-proxy/config"
	"apollo-proxy/model"
	"github.com/philchia/agollo/v4"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	MAX_COUNT = 3600
)

var (
	/**
	第一次启动的时候，所有配置请求都会触发update回调，为了防止无意义的任务创建，等指定时间后触发的update回调，再创建jenkins任务
	*/
	//initDuration  = 1 * time.Minute
	initDuration  = 10 * time.Second
	isAllowUpdate = false
	lock          sync.Mutex
	apps          map[string]bool

	goroutines = 5
	chans      chan struct{}
)

/**
构建工具要实现的接口
*/
type BuildTool interface {
	Build(*model.Project) (int64, error)
	IsFinish(jobId int64) (bool, error)
	JobStatus() string
	Result() string
}

func init() {
	apps = make(map[string]bool)

	chans = make(chan struct{}, goroutines)

}

/**
apollo 配置更新通知
*/
func Watch() {
	go func() {
		WatchApollo()
		/*
			5分钟重新加载一次数据
		*/
		ticker := time.NewTicker(5 * time.Minute)
		for {
			select {
			case <-ticker.C:
				WatchApollo()
			}
		}
	}()

}

/**
获取应用信息，并开始监听变化
*/
func WatchApollo() {
	allApp := model.JenkinsApps()
	lock.Lock()
	defer lock.Unlock()
	for _, app := range allApp {
		cacheKey := app.Project + ":" + app.Env
		if _, ok := apps[cacheKey]; ok {
			continue
		}
		apps[cacheKey] = true
		WatchProjects(app.FindProjects())
	}
	time.AfterFunc(initDuration, func() {
		log.Println("allow update")
		isAllowUpdate = true
	})
}

func WatchProjects(projects []*model.Project) {
	if projects == nil {
		return
	}
	project := projects[0]
	conf := &agollo.Conf{
		MetaAddr: config.ApolloMap[strings.ToLower(project.Env)],
		CacheDir: project.Env,
		AppID:    project.ApolloAppId,
		Cluster:  project.ClusterName,
	}

	namespace := make([]string, 0)
	for _, project := range projects {
		namespace = append(namespace, project.Namespace)
	}
	conf.NameSpaceNames = namespace
	client := agollo.NewClient(conf)
	client.OnUpdate(func(event *agollo.ChangeEvent) {
		if !isAllowUpdate {
			return
		}
		log.Println("Create Job")
		// 创建构建任务
		CreateBuild(event.Namespace, project)
	})
	client.Start()
}

func CreateBuild(namespace string, project *model.Project) {
	buildTools := GetBuildTools()
	for _, build := range buildTools {
		_build := build
		chans <- struct{}{}
		go func() {
			defer func() {
				<-chans
			}()
			buildId, err := _build.Build(project)
			if buildId == 0 {
				// 项目不存在
				if err != nil {
					log.Println("create build fail: ", err)
				}
				return
			}
			buildInfo := NewBuildInfo(project, namespace, buildId)
			for {
				time.Sleep(time.Second)
				finished, err := _build.IsFinish(buildId)
				if err != nil {
					log.Println("err: ", err)
					break
				}
				buildInfo.State = _build.JobStatus()
				buildInfo.Result = _build.Result()
				buildInfo.Count += 1
				buildInfo.Update()
				if finished || buildInfo.Count >= MAX_COUNT {
					break
				}
			}
		}()
	}
}

func GetBuildTools() []BuildTool {
	buildTools := make([]BuildTool, 0)
	jenkinsInstance := NewJenkins()
	if jenkinsInstance != nil {
		buildTools = append(buildTools, jenkinsInstance)
	}
	gitlabInstance := NewGitLab()
	if gitlabInstance != nil {
		buildTools = append(buildTools, gitlabInstance)
	}
	return buildTools
}

func NewBuildInfo(project *model.Project, namespace string, buildId int64) model.BuildInfo {
	buildInfo := model.BuildInfo{
		Project:   project.Project,
		Namespace: namespace,
		Env:       project.Env,
		Branch:    project.Branch,
		JobId:     int(buildId),
		State:     "",
	}
	buildInfo.Save()
	return buildInfo
}
