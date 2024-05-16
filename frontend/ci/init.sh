#!/usr/bin/env bash
# 启用 POSIX 模式并设置严格的错误处理机制
set -o posix errexit -o pipefail

NAME="frontend"
NAMESPACE="frontend"
PORT1="80"
PORT2="443"
IMAGE="ccr.ccs.tencentyun.com/lisa/frontend:v2"
PORT_TYPE="LoadBalancer"

cat > namespace.yml <<EOF
# 命名空间
apiVersion: v1
kind: Namespace
metadata:
  name: ${NAMESPACE}
EOF

cat > service.yml <<EOF
# 服务清单
apiVersion: v1
kind: Service
metadata:
  name: ${NAMESPACE}-service
  namespace: ${NAMESPACE}
spec:
  type: ${PORT_TYPE}
  ports:
    - port: ${PORT1}
      targetPort: ${PORT1}
      protocol: TCP
      name: http
    - port: ${PORT2}
      targetPort: ${PORT2}
      protocol: TCP
      name: https
  selector:
    app: ${NAME}
EOF

cat > deployment.yml <<EOF
# 部署清单
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${NAME}-deployment
  namespace: ${NAMESPACE}
  labels:
    app: ${NAME}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx-quic
          image: ${IMAGE}  # 替换成您的镜像
          ports:
            - containerPort: ${PORT1}
            - containerPort: ${PORT2}
          volumeMounts:
            - name: html-volume
              mountPath: /etc/nginx/html
            - name: ssl-volume
              mountPath: /etc/nginx/ssl
            - name: conf-volume
              mountPath: /etc/nginx/conf.d
      volumes:
        - name: html-volume
          hostPath:
            path: /home/nginx/html
        - name: ssl-volume
          hostPath:
            path: /home/nginx/ssl
        - name: conf-volume
          hostPath:
            path: /home/nginx/conf
EOF
