#!/usr/bin/bash

set -x
# login
# argocd login 192.168.2.155:30618 --insecure --grpc-web

export PROJECT_NAME="microservice-shop"
export PROJECT_REPO="http://192.168.2.158:7080/root/test.git"

kubectl create ns $PROJECT_NAME

argocd app create $PROJECT_NAME \
  --repo $PROJECT_REPO \
  --path . \
  --dest-server https://kubernetes.default.svc \
  --dest-namespace $PROJECT_NAME

# 与argo同步一次
argocd app sync $PROJECT_NAME

# 自动同步, 默认每三分钟检查git repo
argocd app set $PROJECT_NAME --sync-policy automated

# 列出
argocd app list

# 获取
kubectl get all -n $PROJECT_NAME
echo 111

set +x
