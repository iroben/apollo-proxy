package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type AppConfig struct {
	App struct {
		Port   int
		Mode   string
		Token  string
		User   string
		Passwd string
	}

	Jenkins struct {
		User         string
		Token        string
		Domain       string
		TriggerKey   string `yaml:"trigger_key"`
		TriggerValue string `yaml:"trigger_value"`
	}
	GitLab struct {
		Token        string
		Domain       string
		TriggerKey   string `yaml:"trigger_key"`
		TriggerValue string `yaml:"trigger_value"`
	}
	Apollo struct {
		Dev  string
		Fat  string
		Prod string
	}
	Mysql struct {
		Host   string
		User   string
		Passwd string
		DbName string
	}
}

var ApolloMap = make(map[string]string)

var Config AppConfig

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
func init() {
	InitConfig()
}

func InitConfig() {
	configPath := []string{".", "./conf", "../conf"}
	for _, path := range configPath {
		configFile := path + "/app.yaml"
		if Exist(configFile) {
			loadConfig(configFile)
			return
		}
	}
	log.Println("配置初始化失败")
}

func loadConfig(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Println("配置文件不存在")
		return
	}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&Config); err != nil {
		log.Println("配置解析失败：", err)
		return
	}
	ApolloMap["dev"] = Config.Apollo.Dev
	ApolloMap["fat"] = Config.Apollo.Fat
	ApolloMap["prod"] = Config.Apollo.Prod
}
