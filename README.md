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
    ...
    "bundle": {
      ...
      "identifier": "com.tauri.dev",
      ...
    }
    ...
}
```
修改后的值:
```json
{
    ...
    "bundle": {
      ...
      "identifier": "com.tauri.build",
      ...
    }
    ...
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
构建 Web应用为 Docker 镜像
编辑
> 将`image`: `web`  # 这里需要替换成你使用 `Docker build`的名字

语法:
```shell
docker build [--progress=plain] [--no-cache] -t <image_name> <path>
```

例: 在当前目录构建所有文件, 镜像名为`web`, 不使用`Docker`缓存, 显示构建过程的详细信息
- --progress=plain: 构建过程中显示的详细信息的格式
- --no-cache : 不使用缓存
- web : 构建的镜像名
- . : 当前目录所有文件
```shell
docker build --progress=plain --no-cache -t web .
```

- -f: 指定文件
- -d: 后台执行
- --build: 强制重新构建

```shell
docker-compose up \
-f docker-compose.yml \
-d \
--build
```

#### 后端
在项目的根目录执行以下脚本

>  如果项目共用一个 go.mod 模块文件, 能够得到完好支持
> 
>  `Docker build` 执行时的上下文环境是依据当前执行的目录为主目录, 并非 `Dockerfile` 的目录, 例如当前目录为`/app`,  `Dockerfile` 所在的目录是`/app/user/`,执行时是使用`/app`作为主目录, 而不是 `Dockerfile` 的目录
> 
> Docker的新版已经删除--progress=plain, 如果遇到--progress=plain的错误, 删掉即可

```shell
docker build -f <Dockerfile-path> --progress=plain --no-cache -t <image-name> .
```
修改`Dockerfile`的`ENTRYPOINT`, 把`<appname>` 替换为go mod xxx生成的名称, 也就是go build 之后的应用名称, 例如`backend`

示例: 
修改前:
```Dockerfile
command: ["/app/<appname>", "-conf", "/data/conf"]
```

修改后:
```Dockerfile
command: ["/app/backend", "-conf", "/data/conf"]
```

## CI/CD
本项目的CI是分离结构, 如果你是前后端开发, 可以把CI文件统一到根目录

### Gitlab CI


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
