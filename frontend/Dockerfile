# 要依赖的node版本, 这里以node 16演示, 使用最小的node包alpine版本减少镜像大小
FROM node:16-alpine AS builder

WORKDIR /web
COPY . /web

# 分离环境
RUN cd /web

# 可选, 如果你的服务器环境连npm都下不动依赖包, 尝试使用大陆proxy
RUN npm config set registry https://registry.npmmirror.com

# 本项目使用 pnpm 管理包, 根据你的实际来修改
RUN npm install pnpm -g
# 验证 pnpm安装与否
RUN pnpm -v

# 可选, 如果你的服务器环境连pnpm都下不动依赖包, 尝试使用大陆proxy
#RUN pnpm config set registry https://registry.npmmirror.com

# 删除 npm 下载的包,由 pnpm 接管
RUN rm -f node_modules
# 安装依赖
RUN pnpm install
# RUN pnpm install --no-frozen-lockfile 不根据package.lock的包来下载

# 处理下载出错
RUN RET=$? if [[ RET -ne 0 ]]; then tail -n 100 $HOME/.npm/_logs/*.log; exit $RST; fi;

# 打包
RUN pnpm build

# 使用nginx代理 web 项目的映射
FROM nginx:latest

# 把上一步骤打包好的的 dist目录传递到 nginx 默认的 html 目录作为映射
COPY --from=builder /web/dist/ /usr/share/nginx/html/

# 挂载 Nginx 的必要文件,
# /etc/nginx/ssl/是 SSL 配置,不需要删掉即可
# /var/log/nginx/ 是 nginx 日志,不需要删掉即可
# /var/cache/nginx/是 nginx 缓存, 不需要删掉即可
# /etc/nginx/conf/ 是 nginx 的配置文件, 必须保留
# /usr/share/nginx/html/是 nginx 的 html 目录, 必须保留
VOLUME ["/var/cache/nginx/", "/etc/nginx/ssl/", "/etc/nginx/conf/", "/etc/nginx/conf/nginx.conf", "/var/log/nginx/", "/usr/share/nginx/html/"]

# 运行 nginx 服务
CMD ["nginx", "-g", "daemon off;"]
