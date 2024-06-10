#!/usr/bin/env bash
# 启用 POSIX 模式并设置严格的错误处理机制
set -o posix errexit -o pipefail

# https://www.jaegertracing.io/docs/1.57/deployment/#all-in-one
## make sure to expose only the ports you use in your deployment scenario!
if command -v docker-compose > /dev/null 2>&1; then
    echo "docker-compose is installed."
    cat > compose.yml <<EOF
services:
  jaeger:
    image: jaegertracing/all-in-one:latest
    container_name: jaeger
    environment:
      - COLLECTOR_OTLP_ENABLED=true
      - COLLECTOR_ZIPKIN_HOST_PORT=:9411
    ports:
      - "0.0.0.0:5775:5775/udp"
      - "0.0.0.0:6831:6831/udp"
      - "0.0.0.0:6832:6832/udp"
      - "0.0.0.0:5778:5778"
      - "0.0.0.0:16686:16686"
      - "0.0.0.0:14250:14250"
      - "0.0.0.0:14268:14268"
      - "0.0.0.0:14269:14269"
      - "0.0.0.0:4317:4317"
      - "0.0.0.0:4318:4318"
      - "0.0.0.0:9411:9411"
EOF
    docker-compose -f compose.yml down
    docker-compose -f compose.yml up -d
    docker-compose -f compose.yml logs
else
    echo "docker-compose is not installed."
    docker stop jaeger
    docker rm jaeger
    docker run -itd \
      --name jaeger \
      -e COLLECTOR_OTLP_ENABLED=true \
      -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
      -p "0.0.0.0:5775:5775/udp" \
      -p "0.0.0.0:6831:6831/udp" \
      -p "0.0.0.0:6832:6832/udp" \
      -p "0.0.0.0:5778:5778" \
      -p "0.0.0.0:16686:16686" \
      -p "0.0.0.0:14250:14250" \
      -p "0.0.0.0:14268:14268" \
      -p "0.0.0.0:14269:14269" \
      -p "0.0.0.0:4317:4317" \
      -p "0.0.0.0:4318:4318" \
      -p "0.0.0.0:9411:9411" \
      jaegertracing/all-in-one:latest
    docker ps
    docker logs -f jaeger
fi
