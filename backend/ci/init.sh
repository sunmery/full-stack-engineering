#!/usr/bin/env bash
# 启用 POSIX 模式并设置严格的错误处理机制
set -o posix errexit -o pipefail

NAME="backend"
NAMESPACE="backend"
PORT1="30001"
PORT2="30002"
IMAGE="ccr.ccs.tencentyun.com/lisa/go:full"
PORT_TYPE="LoadBalancer"
LIVENESS_PROBE_PATH="/helloworld/lisa"

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
          # Liveness Probe配置
          # 目的是定期检查容器内应用的健康状况。如果探测失败（达到failureThreshold设定的次数），
          # Kubernetes将认为容器不健康并自动重启该容器。
          livenessProbe:
            # 通过HTTP GET请求检查
            httpGet:
              # 探测的路径，应指向一个能快速响应且表明应用运行正常的端点
              path: ${LIVENESS_PROBE_PATH}/lisa
              # 应用监听的端口
              port: ${PORT1}
            # 部署后首次探测前等待的时间，给予应用足够的启动时间
            initialDelaySeconds: 60
            # 探测间隔时间，每隔多久进行一次检查
            periodSeconds: 60
            # 请求超时时间，超过此时间认为探测失败
            timeoutSeconds: 5
            # 探测失败次数的阈值，连续失败达到此次数后，将采取相应动作（此处为重启容器）
            failureThreshold: 3

          # Readiness Probe配置
          # 用于判断容器是否已准备好接收外部请求。如果未准备好（探测失败），Kubernetes不会将流量路由到该容器，直到探测成功。
          readinessProbe:
            # 通过HTTP GET方式检查
            httpGet:
              # 确保应用能够处理即将到来的请求
              path: ${LIVENESS_PROBE_PATH}
              port: ${PORT1}
            # 相较于Liveness Probe，Readiness Probe可以更快开始，因为只需等待应用启动的基本就绪
            initialDelaySeconds: 15
            # 更频繁地检查，以便迅速响应容器就绪状态的变化
            periodSeconds: 5
            # 缩短超时时间，加快反馈速度
            timeoutSeconds: 3
            # 连续失败达到此次数，Kubernetes将认为容器尚未准备好服务请求
            failureThreshold: 3
          volumeMounts:
            - name: config-volume
              mountPath: /data/conf
          command: [ "/app/${NAME}", "-conf", "/data/conf" ]
      volumes:
        - name: config-volume
          hostPath:
            path: /path/to/configs
EOF
