package service

import (
	"apollo-proxy/config"
	"apollo-proxy/model"
	"encoding/json"
	"github.com/xanzy/go-gitlab"
	"log"
)

type GitLab struct {
	Instance *gitlab.Client
	project  *model.Project
	pipeline *gitlab.Pipeline
}

func NewGitLab() *GitLab {
	gitlabConfig := config.Config.GitLab
	if gitlabConfig.Domain == "" {
		return nil
	}
	instance, err := gitlab.NewClient(gitlabConfig.Token, gitlab.WithBaseURL(gitlabConfig.Domain))
	if err != nil {
		log.Fatal("Failed to create client: ", err)
		return nil
	}
	return &GitLab{
		Instance: instance,
	}
}

func (j *GitLab) Build(project *model.Project) (int64, error) {
	pipeline, resp, err := j.Instance.Pipelines.CreatePipeline(project.Project, &gitlab.CreatePipelineOptions{
		Ref: &project.Branch,
		Variables: []*gitlab.PipelineVariable{
			{
				Key:          config.Config.GitLab.TriggerKey,
				Value:        config.Config.GitLab.TriggerValue,
				VariableType: string(gitlab.EnvVariableType),
			}},
	})
	if err != nil {
		log.Println("gitlab create pipeline fail: ", err)
		return 0, err
	}
	if resp.StatusCode == 404 {
		return 0, nil
	}
	j.project = project
	return int64(pipeline.ID), nil
}

func (j *GitLab) IsFinish(pipelineId int64) (bool, error) {
	pipeline, _, err := j.Instance.Pipelines.GetPipeline(j.project.Project, int(pipelineId))
	if err != nil {
		log.Println("gitlab get pipeline fail: ", err)
		return false, err
	}
	j.pipeline = pipeline
	return pipeline.FinishedAt != nil, nil
}

func (j *GitLab) JobStatus() string {
	return j.pipeline.Status
}

func (j *GitLab) Result() string {
	bt, _ := json.MarshalIndent(j.pipeline, "", "    ")
	return string(bt)
}
