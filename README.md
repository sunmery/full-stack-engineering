## 介绍
### 前端
前端技术栈使用React作为前端kit, 考量有: React灵活, 前端的生态最丰富
- 桌面库: 采用tauri, 一个rust打造的桌面, 占用极小
- CLI: 采用Vite, 现代和快速

### 后端
- 语言选择: 使用Golang的基于DDD思想的Kratos作为后端项目的kit, 即可作为单体架构使用, 也可以作为微服务架构使用
- 协议: HTTP/RPC, 支持HTTP与protobuf协议

### 运维
#### CI
支持Gitlab, Github的CI, 内置ci脚本, 根据实际需求修改

#### CD
本项目使用Gitlab + Argo 作为自动化的部署

#### 网关与安全
TODO

## 前端工程化
参考[文章](https://juejin.cn/post/7346610589987520512)

## 后端工程化
参考[文章](https://juejin.cn/post/7346519210724278306)

## 运维工程化
参考[文章](https://juejin.cn/post/7346610589987569664)

## 运行Web项目
1. 初始化git仓库, 不执行husky会报错
    ```shell
    git init
    ```

2. 安装依赖
    ```shell
    pnpm i
    ```

- 本地运行
```shell
pnpm dev
```

- 运行Web桌面应用

> 构建桌面应用需要Rust环境

```shell
pnpm tauri dev
```

- 打包
```shell
pnpm build
```

- 构建Web桌面应用

> 构建桌面应用需要Rust环境

编辑`project/frontend/src-tauri/tauri.conf.json`文件的tauri.conf.json > tauri > bundle > identifier的"com.tauri.dev"值改为:""com.tauri.build""
原来的值:
```json
{
    
    "bundle": {
      
      "identifier": "com.tauri.dev",
      
    }
    
}
```
修改后的值:
```json
{
    
    "bundle": {
      
      "identifier": "com.tauri.build",
      
    }
    
}
```
运行构建
```shell
pnpm tauri build
```

## 运行后端项目
### 准备环境
- Golang
- https://go-kratos.dev/docs/getting-started/start

1. 安装依赖
```shell
kratos upgrade
```

- 运行
```shell
kratos run
```

### 创建微服务
- <path>/<project>: 创建的路径, 例如`service/account`
- -r <https://xxx.git>: 指定的`kratos-layout`的git仓库URL, 例如gitee的`kratos-layout`的URL: `https://gitee.com/go-kratos/kratos-layout.git`
- --nomod: 共用 go.mod, 大仓模式
- -b <brach>: 可选, 指定分支, 例如`main`
```shell
kratos new <path>/<project> [-r <https://xxx.git>] [-b <brach>] [--nomod]
```

## 部署

[运维工程化](https://juejin.cn/post/7346610589987569664)有详细说明, 这里是简化说明

### Kubernetes
#### 前端
修改`deploy.yml`,把`spec.template.spec.containers.image`字段替换成你的镜像URL地址, 例如`https://example.com/v2/web`

#### 后端
> Docker的新版已经删除--progress=plain, 如果遇到--progress=plain的错误, 删掉即可

注意事项:
- 请将 /path/to/configs 替换为实际的本地配置文件路径。
- 将 <appname> 替换为实际的应用程序名称。
- 根据需要调整 replica 数量、Service 类型等配置

最后在安装了Kubernetes的服务器环境上执行:
```shell
kubectl apply -f backend/deployment.yaml
```

### Docker
#### 前端
##### 构建镜像
编辑
> 将`image`: `web`  # 这里需要替换成你使用 `Docker build`的名字

语法:
```shell
docker build [--progress=plain] [--no-cache] -t <user>/<image_name>:<tag> <path>
```

例: 在当前目录构建所有文件, 镜像名为`web`, 不使用`Docker`缓存, 显示构建过程的详细信息
- --progress=plain: 以简洁的方式显示构建进度信息
- --no-cache : 不使用缓存
- web : 构建的镜像名
- . : 当前目录所有文件
```shell
docker build --progress=plain --no-cache -t web .
```

##### 运行镜像
###### docker-compose
选项: 
- -f: 指定文件
- -d: 后台执行
- --build: 强制重新构建

```shell
docker-compose up \
-f docker-compose.yml \
-d \
--build
```

###### docker run
选项:
-d：以后台模式运行容器。
--name backend：指定容器的名称为 backend。
--restart unless-stopped：定义容器的重启策略为除非手动停止，否则总是重新启动。
-p 80:80 -p 443:443：将容器内部的 80 端口映射到宿主机的 80 端口，将容器内部的 443 端口映射到宿主机的 443 端口。
mandala/frontend:v1：指定要运行的镜像为 mandala/frontend 的版本 v1。

```shell
docker run -d \
    --name backend \
    --restart unless-stopped \
    -p 80:80 \
    -p 443:443 \
    mandala/frontend:v1
```

#### 后端
##### 构建镜像
在项目的根目录执行以下脚本

>  如果项目共用一个 go.mod 模块文件, 能够得到完好支持
> 
>  `Docker build` 执行时的上下文环境是依据当前执行的目录为主目录, 并非 `Dockerfile` 的目录, 例如当前目录为`/app`,  `Dockerfile` 所在的目录是`/app/user/`,执行时是使用`/app`作为主目录, 而不是 `Dockerfile` 的目录
> 
> Docker的新版已经删除--progress=plain, 如果遇到--progress=plain的错误, 删掉即可

语法:
```shell
docker build [--progress=plain] [--no-cache] -t <user>/<image_name>:<tag> <path>
```

- --progress=plain: 构建过程中显示的详细信息的格式
- --no-cache: 不使用缓存
- -t: 标签, 例如 myusername/myimage:v1
- web : 构建的镜像名
- . : 当前目录所有文件
示例:
```shell
docker build -f Dockerfile --progress=plain --no-cache -t mandala/backend:v1 .
```
修改`Dockerfile`的`ENTRYPOINT`, 把`backend` 替换为go mod xxx生成的名称, 也就是go build 之后的应用名称, 例如`backend`

示例:
```Dockerfile
command: ["/app/backend", "-conf", "/data/conf"]
```

##### 运行镜像

###### docker-compose
修改`docker-compose.yml`:
1. 修改`image`:字段为你的镜像URL, 如果已经打包上传到远程容器注册表, 这里就是镜像的完整URL地址.
如果没有, 那么你需要先使用`docker buiuld`进行操作, 然后填写你在`docker buiuld`时填写的标签, 例如上一步操作的`mandala/backend:v1`

进入到`docker-compose.yml`所在的目录, 然后运行该应用的容器:

> 建议使用docker-compose的方式, 这种方便记录你的操作, 可以存留你的操作, 而且也方便其他人员查看你的部署文件

docker-compose:
```shell
docker-compose up -d
```

###### docker run
```shell
docker run -d \
    --name backend \
    --restart unless-stopped \
    -p 30001:30001 \
    -p 30002:30002 \
    -v $(pwd)/configs:/data/conf \
    mandala/backend:v1 /app/backend -conf /data/conf
```

## CI/CD
本项目的CI是分离结构, 如果你是前后端开发, 可以把CI文件统一到根目录

### Gitlab CI

#### 准备环境
在你的服务器下安装, 在我的[运维工程化](https://juejin.cn/post/7346610589987569664)文章里面详细写了如何安装, 这里不在赘述
- Gitlab
- Gitlab Runner

1. 定义变量

Gitlab -> 项目页 -> 设置 -> CI/CD -> 变量 -> 添加:
- REGISTER_ADDRESS: 注册表的地址, 例如`ccr.ccs.tencentyun.com`
- REGISTER_USERNAME: 注册表的用户名
- REGISTER_PASSWORD: 注册表的密码
- REGISTER_REPO: 应用在注册表的空间名, 必须在注册表提前创建, 例如`web`
- SSHPASS: 远程服务器的密码, 是你要上传文件的服务器的密码
![img.png](gitlab_var.png)

2. 编辑前端/后端项目目录下的`.gitlab-ci.yml`
把`variables`里面的变量值都改成你实际的变量值

### Github CI
#### 前端
根据你的实际需求, 修改CI文件
```shell
frontend/.github/workflows/ci.yaml
```

#### 后端
根据你的实际需求, 修改CI文件
```shell
backend/.github/workflows/ci.yaml
```
