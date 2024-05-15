#!/usr/bin/env bash
# 启用 POSIX 模式并设置严格的错误处理机制
set -o posix errexit -o pipefail

NAME="backend"
NAMESPACE="backend"
PORT1="30001"
PORT2="30002"
IMAGE="ccr.ccs.tencentyun.com/lisa/go:full"
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
  name: ${NAME}-service
spec:
  selector:
    app: ${NAME}
  ports:
    - protocol: TCP
      port: ${PORT1}
      targetPort: ${PORT1}
      nodePort: ${PORT1}
    - protocol: TCP
      port: ${PORT2}
      targetPort: ${PORT2}
      nodePort: ${PORT2}
  type: ${PORT_TYPE}
EOF

cat > deployment.yml <<EOF
# 部署清单
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${NAME}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ${NAME}
  template:
    metadata:
      labels:
        app: ${NAME}
    spec:
      containers:
        - name: ${NAME}
          image: ${IMAGE}
          ports:
            - containerPort: ${PORT1}
            - containerPort: ${PORT2}
          volumeMounts:
            - name: config-volume
              mountPath: /data/conf
          command: [ "/app/${NAME}", "-conf", "/data/conf" ]
      volumes:
        - name: config-volume
          hostPath:
            path: /path/to/configs
EOF
