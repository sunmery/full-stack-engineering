#!/usr/bin/env bash

# 启用 POSIX 模式并设置严格的错误处理机制
set -o posix errexit -o pipefail

unset ARGOCD_NAMESPACE
unset ROLE_NAME
unset BRANCH
unset PROJECT_GIT_URL

# argocd所在的的命名空间
export ARGOCD_NAMESPACE="argo"
# 角色名称, 用于管理项目
export ROLE_NAME="admin"
# Kubernetes集群地址
export CLUSTER_SERVER="https://192.168.2.160:6443"
#export CLUSTER_SERVER="https://kubernetes.default.svc"

# 仓库URL
export BRANCH="main"
# 仓库地址
export PROJECT_GIT_URL="https://gitlab.com/lookeke/full-stack-engineering.git"

# 前端命名空间, 不需要额外创建命名空间选择default即可
export FRONTEND_NAMESPACE="frontend"
# argocd中的前端项目名, 用于分配团队人员的操作权限
export FRONTEND_PROJECT_NAME="frontend"
# 前端应用的名称
export FRONTEND_APPLICATION_NAME="react"
# Kubernetes 资源清单在仓库中的路径, 相对于仓库根目录的路径
export FRONTEND_DEPLOY_PATH="frontend/ci"

# 后端端命名空间, 不需要额外创建命名空间选择default即可
export BACKEND_NAMESPACE="backend"
# argocd中的后端项目名, 用于分配团队人员的操作权限
export BACKEND_PROJECT_NAME="backend"
# 后端应用的名称
export BACKEND_APPLICATION_NAME="go"

# 获取Git Repo URL
if [ -z "$BRANCH" ]; then
  echo "用户未设置Git Repo URL, 尝试自动获取"
  # 该项目所在的git仓库地址
  export BRANCH=$(git rev-parse --abbrev-ref HEAD)
fi

if [ -z "$PROJECT_GIT_URL" ]; then
  if command -v git &> /dev/null
  then
      echo "Git is installed."
      export PROJECT_GIT_URL=$(git config --get remote."${BRANCH}".url)
      $PROJECT_GIT_URL
      echo "${PROJECT_GIT_URL}"
  else
      echo "Git is not installed."
      exit 1
  fi
fi

# 创建argocd的Project(项目)的Role(角色)
cat > project-role.yml <<EOF
#  创建角色
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-cm
  namespace: ${ARGOCD_NAMESPACE}
data:
  accounts.${ROLE_NAME}: "apiKey, login"
# kubectl apply -f argocd-cm.yaml -n ${ARGOCD_NAMESPACE}
EOF

# 给Role分配Project的权限
cat > project-rbac.yml <<EOF
# 分配角色给 frontend-group 前端组 和 backend-group 后端组
# 并具有适当的权限
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-rbac-cm
  namespace: ${ARGOCD_NAMESPACE}
data:
  policy.csv: |
    p, role:admin, applications, *, *, allow
    p, role:${ROLE_NAME}, applications, *, *, allow
    g, admin, role:admin
    g, ${ROLE_NAME}, proj:frontend:${ROLE_NAME}
    g, ${ROLE_NAME}, proj:backend:${ROLE_NAME}
# kubectl apply -f argocd-rbac-cm -n ${ARGOCD_NAMESPACE}
EOF

cat > create-frontend-proj.yml <<EOF
apiVersion: argoproj.io/v1alpha1
kind: AppProject
metadata:
  name: ${FRONTEND_NAMESPACE} # 项目名称
  namespace: ${ARGOCD_NAMESPACE} # 项目所在的命名空间，默认为argocd，可根据实际情况调整
spec:
  # 项目描述
  description: "This project is for managing frontend applications"

  # 允许应用部署到的命名空间列表
  destinations:
  - namespace: ${FRONTEND_NAMESPACE}
    server: ${CLUSTER_SERVER} # 集群API地址，示例值需替换为实际地址

  # 源代码仓库配置
  sourceRepos:
  - ${PROJECT_GIT_URL} # 允许使用的Git仓库地址，根据实际情况修改

  # 角色与成员
  roles:
    - name: ${ROLE_NAME} # 角色名称
      policies:
      # 允许admin角色在frontend项目中对所有应用进行所有操作
      - p, proj:${FRONTEND_PROJECT_NAME}:${ROLE_NAME}, applications, *, ${FRONTEND_PROJECT_NAME}/*, allow
      # 允许查看frontend命名空间相关的集群信息
      - p, proj:${FRONTEND_PROJECT_NAME}:${ROLE_NAME}, clusters, get, ${FRONTEND_PROJECT_NAME}/*, allow

  # 可选：项目默认值
  # 这些默认值可以被应用级别的设置覆盖
  # clusterResourceWhitelist, namespaceResourceBlacklist等可以根据需要添加
EOF

cat > create-backend-proj.yml <<EOF
apiVersion: argoproj.io/v1alpha1
kind: AppProject
metadata:
  name: ${BACKEND_NAMESPACE} # 项目名称
  namespace: ${ARGOCD_NAMESPACE} # 项目所在的命名空间，默认为argocd，可根据实际情况调整
spec:
  # 项目描述
  description: "This project is for managing BACKEND applications"

  # 允许应用部署到的命名空间列表
  destinations:
  - namespace: ${BACKEND_NAMESPACE}
    server: ${CLUSTER_SERVER} # 集群API地址，示例值需替换为实际地址

  # 源代码仓库配置
  sourceRepos:
  - ${PROJECT_GIT_URL} # 允许使用的Git仓库地址，根据实际情况修改

  # 角色与成员
  roles:
  - name: ${ROLE_NAME} # 角色名称
    # 注意这里的subjects配置应当符合Argo CD的RBAC规范，例如使用proj:backend:admin
    # 定义角色权限
    policies:
    - p, proj:${BACKEND_PROJECT_NAME}:${ROLE_NAME}, applications, *, ${BACKEND_PROJECT_NAME}/*, allow # 允许admin角色在BACKEND项目中对所有应用进行所有操作
    - p, proj:${BACKEND_PROJECT_NAME}:${ROLE_NAME}, clusters, get, ${BACKEND_PROJECT_NAME}/*, allow # 允许查看集群信息

  # 可选：项目默认值
  # 这些默认值可以被应用级别的设置覆盖
  # clusterResourceWhitelist, namespaceResourceBlacklist等可以根据需要添加
EOF

# 创建前端项目应用
cat > application-frontend.yml <<EOF
apiVersion: argoproj.io/v1alpha1  # 指定 Argo CD API 版本
kind: Application  # 定义资源类型为 Application
metadata: # 元数据部分
  name: ${FRONTEND_APPLICATION_NAME}   # 指定 Application 的名称
  namespace: ${ARGOCD_NAMESPACE}  # argocd所属的命名空间
  # 定义资源的 finalizers
  # https://argo-cd.readthedocs.io/en/stable/user-guide/app_deletion/#about-the-deletion-finalizer
  finalizers:
    - resources-finalizer.argocd.argoproj.io  # 删除时行级联删除
    #- resources-finalizer.argocd.argoproj.io/background  # 删除时后台行级联删除
spec: # 规范部分
  project: ${FRONTEND_PROJECT_NAME}  # 应用程序将被配置的项目名称，这是在 Argo CD 中应用程序的一种组织方式
  source: # 指定源
    # Kubernetes 资源清单在仓库中的路径
    path: ${FRONTEND_DEPLOY_PATH}
    # 指定 Git 仓库的 URL
    repoURL: ${PROJECT_GIT_URL}
    # 使用的 git 分支
    targetRevision: ${BRANCH}
  # 部署应用到Kubernetes 集群中的位置
  destination:
    namespace: ${FRONTEND_NAMESPACE}  # 指定应用的命名空间
    server: ${CLUSTER_SERVER}  # 如果部署到同一集群，可以省略
  syncPolicy: # 指定同步策略
    automated: # 自动化同步
      prune: true  # 启用资源清理
      selfHeal: true  # 启用自愈功能
      allowEmpty: false  # 禁止空资源
    syncOptions: # 同步选项
      - Validate=false  # 是否启用验证
      - CreateNamespace=true  # 启用创建命名空间
    retry: # 重试策略
      limit: 5  # 重试次数上限
      backoff: # 重试间隔
        duration: 10s  # 初始重试间隔
        factor: 2  # 重试间隔因子
        maxDuration: 3m  # 最大重试间隔
EOF

cat > project-app.yml <<EOF
apiVersion: argoproj.io/v1alpha1
kind: AppProject
metadata:
  name: ${FRONTEND_NAMESPACE}
  namespace: ${ARGOCD_NAMESPACE}
spec:
  # 说明
  description: Project for frontend applications
  # 允许的集群
  destinations:
    # 允许的命名空间
    - name: ${FRONTEND_PROJECT_NAME}
      namespace: "${FRONTEND_NAMESPACE}"
      # 集群地址
      server: ${CLUSTER_SERVER}
    - name: ${BACKEND_PROJECT_NAME}
      namespace: "{BACKEND_NAMESPACE}"
      # 集群地址
      server: ${CLUSTER_SERVER}
  # 允许的目标 K8s 资源类型
  clusterResourceWhitelist:
    - group: '*'
      #kind: '*'
      kind: Namespace
  # Allow all namespaced-scoped resources to be created, except for ResourceQuota, LimitRange, NetworkPolicy
  #namespaceResourceBlacklist:
  #  - group: ''
  #    kind: ResourceQuota
  #  - group: ''
  #    kind: LimitRange
  #  - group: ''
  #    kind: NetworkPolicy
  # Deny all namespaced-scoped resources from being created, except for Deployment and StatefulSet
  #namespaceResourceWhitelist:
  #  - group: 'apps'
  #    kind: Deployment
  #  - group: 'apps'
  #    kind: StatefulSet
  sourceRepos:
    - "*"
  roles:
    - name: ${ROLE_NAME}
      description: Access role for ROLE_NAME user
      policies:
        - p, proj:default:${ROLE_NAME}, applications, *, ${FRONTEND_PROJECT_NAME}/*, allow
        - p, proj:${FRONTEND_PROJECT_NAME}:${ROLE_NAME}, applications, get, ${FRONTEND_PROJECT_NAME}/*, allow
        - p, proj:${FRONTEND_PROJECT_NAME}:${ROLE_NAME}, applications, create, ${FRONTEND_PROJECT_NAME}/*, allow
        - p, proj:${FRONTEND_PROJECT_NAME}:${ROLE_NAME}, applications, sync, ${FRONTEND_PROJECT_NAME}/*, allow
        - p, proj:${FRONTEND_PROJECT_NAME}:${ROLE_NAME}, applications, delete, ${FRONTEND_PROJECT_NAME}/*, allow
        - p, proj:${FRONTEND_PROJECT_NAME}:${ROLE_NAME}, repositories, *, ${FRONTEND_PROJECT_NAME}/*, allow
        - p, proj:${FRONTEND_PROJECT_NAME}:${ROLE_NAME}, clusters, get, ${FRONTEND_PROJECT_NAME}/*, allow
  orphanedResources:
    warn: true
# kubectl apply -f project-app.yml
EOF