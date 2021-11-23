package controller

import (
	"apollo-proxy/config"
	"apollo-proxy/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Apollo struct {
	BaseController
}

func (c *Apollo) Configs(ctx *gin.Context) {
	project := new(model.Project)
	if err := ctx.ShouldBindUri(project); err != nil {
		ctx.Status(400)
		return
	}
	project.Save()
	ReverseProxy(config.ApolloMap[strings.ToLower(project.Env)], ctx)
}

type Notification struct {
	AppId         string `form:"appId"`
	Cluster       string `form:"cluster"`
	Notifications string `form:"notifications"`
	JobName       string `form:"job_name"`
}
type NotifyNameSpace struct {
	NamespaceName  string `json:"namespaceName"`
	NotificationId int    `json:"notificationId"`
}

func (c *Apollo) Notifications(ctx *gin.Context) {
	notifyInfo := new(Notification)
	project := new(model.Project)
	if err := ctx.ShouldBindQuery(notifyInfo); err != nil {
		ctx.Status(400)
		return
	}
	if err := ctx.ShouldBindUri(project); err != nil {
		ctx.Status(400)
		return
	}
	notifyNameSpace := make([]NotifyNameSpace, 0)
	if err := json.Unmarshal([]byte(notifyInfo.Notifications), &notifyNameSpace); err != nil {
		ctx.Status(400)
		return
	}
	if len(notifyNameSpace) > 0 {
		for _, namespace := range notifyNameSpace {
			project.Id = 0
			project.ApolloAppId = notifyInfo.AppId
			project.ClusterName = notifyInfo.Cluster
			project.Namespace = namespace.NamespaceName
			project.Save()
		}
	}
	ReverseProxy(config.ApolloMap[strings.ToLower(project.Env)], ctx)
}
func ReverseProxy(host string, ctx *gin.Context) {
	target, _ := url.Parse(host)
	req := ctx.Request
	proxy := httputil.NewSingleHostReverseProxy(target)

	pathArray := strings.Split(req.URL.Path, "/")
	trimPath := strings.Join(pathArray[4:], "/")
	req.URL.Path = singleJoiningSlash(target.Path, trimPath)
	req.URL.Host = target.Host
	req.Host = target.Host
	req.Header.Del("Authorization")

	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
func ReverseProxyBack(host string, ctx *gin.Context) {
	target, _ := url.Parse(host)
	director := func(req *http.Request) {
		req.Header = ctx.Request.Header
		req.URL.Scheme = target.Scheme
		pathArray := strings.Split(req.URL.Path, "/")
		trimPath := ""
		req.Header.Del("Authorization")
		trimPath = strings.Join(pathArray[4:], "/")
		req.URL.Path = singleJoiningSlash(target.Path, trimPath)
		req.URL.Host = target.Host
		req.Host = target.Host
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
