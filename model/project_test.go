package model

import (

	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJenkinsApps(t *testing.T) {
	apps := JenkinsApps()
	assert.NotNil(t, apps)
	for _, app := range apps {
		fmt.Printf("apollo_app_id=%s env=%s \n", app.ApolloAppId, app.Env)
	}
}

func TestProject_FindProjectByJenkinsAppId(t *testing.T) {
	project := Project{
		ApolloAppId: "order-workflow-front",

	}
	apps := project.FindProjectByApolloId()
	assert.NotNil(t, apps)
	for _, app := range apps {
		fmt.Printf("namespace=%s env=%s \n", app.Namespace, app.Env)
	}
}
