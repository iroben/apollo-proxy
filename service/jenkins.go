package service

import (
	"apollo-proxy/config"
	"apollo-proxy/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bndr/gojenkins"
	"log"
	"net/http"
	"regexp"
	"strings"

	"time"
)

type Jenkins struct {
	Instance *gojenkins.Jenkins
	build    *gojenkins.Build
	jobName  string
	buildId  int64
	result   string
}

func NewJenkins() *Jenkins {
	jenkinsConfig := config.Config.Jenkins
	if jenkinsConfig.Domain == "" {
		return nil
	}
	return &Jenkins{
		Instance: gojenkins.CreateJenkins(nil, jenkinsConfig.Domain,
			jenkinsConfig.User, jenkinsConfig.Token),
	}
}
func GetJobName(project *model.Project) string {
	return project.Project + "/job/" + project.Branch
}

type TriggerResponse struct {
	Message string
	Jobs    map[string]struct {
		Triggered bool
		Id        int64
		Url       string
	}
}

func (j *Jenkins) BuildByTrigger(project *model.Project) (int64, error) {
	domain := strings.TrimSuffix(config.Config.Jenkins.Domain, "/")
	jobName := project.Project + "/" + project.Branch
	buildUrl := domain + "/generic-webhook-trigger/invoke" + "?token=" + jobName
	resp, err := http.Post(buildUrl, "Content-Type: application/json", strings.NewReader(
		fmt.Sprintf(`{"%s":"%s"}`,
			config.Config.Jenkins.TriggerKey, config.Config.Jenkins.TriggerValue)))
	if err != nil {
		log.Println("jenkins create build fail: ", err)
		return 0, err
	}
	defer resp.Body.Close()
	var trigger TriggerResponse
	decode := json.NewDecoder(resp.Body)
	if err := decode.Decode(&trigger); err != nil {
		return 0, err
	}
	if buildInfo, ok := trigger.Jobs[jobName]; ok {
		return buildInfo.Id, errors.New("job not exists: " + jobName)
	}
	return 0, nil

}
func (j *Jenkins) BuildJob(project *model.Project) (int64, error) {
	queueId, err := j.Instance.BuildJob(context.TODO(), GetJobName(project), nil)
	if err != nil {
		log.Println("jenkins create build fail: ", err)
		return 0, err
	}
	return queueId, nil

}
func (j *Jenkins) GetJobNameAndBuildId(ctx context.Context, queueId int64) (string, int64) {
	task, err := j.Instance.GetQueueItem(ctx, queueId)
	if err != nil {
		log.Println("jenkins sync build fail: ", err)
		return "", 0
	}
	for task.Raw.Executable.Number == 0 {
		time.Sleep(1000 * time.Millisecond)
		_, err = task.Poll(ctx)
		if err != nil {
			log.Println("jenkins sync build fail: ", err)
			return "", 0
		}
	}
	reg := regexp.MustCompile("job(.*$)")
	data := reg.FindAllStringSubmatch(task.Raw.Task.URL, -1)
	if len(data) != 1 || len(data[0]) != 2 {
		log.Println("jenkins task url match fail: ", err)
		return "", 0
	}
	return data[0][1], task.Raw.Executable.Number
}

func (j *Jenkins) Build(project *model.Project) (int64, error) {
	if config.Config.Jenkins.TriggerKey != "" {
		return j.BuildByTrigger(project)
	}
	return j.BuildJob(project)
}

func (j *Jenkins) IsFinish(queueId int64) (bool, error) {
	ctx := context.TODO()
	if j.jobName == "" {
		jobName, buildId := j.GetJobNameAndBuildId(ctx, queueId)
		if jobName == "" {
			return false, errors.New("get job info fail")
		}
		j.jobName = jobName
		j.buildId = buildId
	}

	build, err := j.Instance.GetBuild(ctx, j.jobName, j.buildId)

	if err != nil {
		log.Println("jenkins get build fail: ", err)
		return false, err
	}
	j.build = build
	return build.GetDuration() > 0, nil
}

func (j *Jenkins) JobStatus() string {
	return j.build.Raw.Result
}

func (j *Jenkins) Result() string {
	bt, _ := json.MarshalIndent(j.build.Raw, "", "    ")
	return string(bt)
}
