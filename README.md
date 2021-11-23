## 为什么要写这个程序
> 1. 接入权限验证功能
> 2. 也是最主要的，我想知道项目和配置的依赖关系，当相关配置更新后，能自动重新构建项目

------

## 快速开始
> 1. clone 项目到本地
> 2. 安装依赖 `go mod vendor`
> 3. 创建好数据库，并修改 **conf/app.yaml** 文件
> 4. 生成数据库信息 `go run cmd/migrate.go`

## 持续集成

`apollo-proxy` 实现了 `jenkins` 和 `gitlab` 的流水线功能

### 集成 `jenkins` 设置（目前只支持**多分支流水线**）
 
1. 生成 `api token`
2. 修改`conf/app.yaml` 的 `jenkins` 信息
3. 如果没有安装 `Generic Webhook Trigger Plugin`，请删除下面两行
    ```
    trigger_key: TRIGGER   
    trigger_value: apollo-proxy
   ```

![jenkins token](https://raw.githubusercontent.com/iroben/apollo-proxy/master/picture/jenkins-token.png)

------

### 集成 `gitlab` 设置

1. 生成 `access token`
2. 修改`conf/app.yaml` 的 `gitlab` 信息

![jenkins token](https://raw.githubusercontent.com/iroben/apollo-proxy/master/picture/gitlab-token.png)

------

## 注意
项目名为 `norecord` 或者分支名为 `dev` 的项目数据，`apollo-proxy` 不会保存


## 记录配置（命名空间）和项目的依赖关系 
![jenkins token](https://raw.githubusercontent.com/iroben/apollo-proxy/master/picture/mysql.png)

## 例子

[前端接入阿波罗配置服务](https://github.com/iroben/apollo-proxy-front-example)

[后端接入阿波罗配置服务](https://github.com/iroben/apollo-proxy-go-example)