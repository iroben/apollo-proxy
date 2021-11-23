package service

import (
	"apollo-proxy/model"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGitLab_Build(t *testing.T) {
	gitlab := NewGitLab()
	buildId, _ := gitlab.Build(&model.Project{
		Project: "667",
		Branch:  "master",
	})
	t.Log("buildId: ", buildId)
	for {
		flag, err := gitlab.IsFinish(buildId)
		if err != nil {
			t.Log("err: ", err)
			break
		}

		fmt.Println("result: ", gitlab.JobStatus())
		time.Sleep(time.Second)
		if flag {
			break
		}
	}
	assert.NotEqual(t, buildId, 0)
}
