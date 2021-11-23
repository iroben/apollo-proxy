package service

import (
	"apollo-proxy/model"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJenkins_Build(t *testing.T) {
	jenkins := NewJenkins()
	buildId, _ := jenkins.Build(&model.Project{
		Project: "front-jenkins-example",
		Branch:  "master",
	})
	t.Log("buildId: ", buildId)
	for {
		flag, err := jenkins.IsFinish(buildId)
		if err != nil {
			t.Log("err: ", err)
			break
		}

		fmt.Println("result: ", jenkins.JobStatus())
		time.Sleep(time.Second)
		if !flag {
			break
		}
	}
	assert.NotEqual(t, buildId, 0)
}
