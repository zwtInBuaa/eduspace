# 基础镜像
FROM golang:1.20-alpine AS builder
# 工作目录
RUN mkdir /build
WORKDIR /build
# 启用CGO
ENV CGO_ENABLED=1
# 换源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    apk add gcc musl-dev && \
    go env -w GOPROXY=https://goproxy.cn,direct
    
COPY go.mod go.sum .
RUN go mod download
# 将当前目录中的所有文件复制到容器的 WORKDIR 目录中
COPY . .
# 构建应用程序
RUN go build -o /build/main .

# 运行时镜像
FROM raspstudio/python-ml:oldxgboost AS runtime

RUN mkdir /app
WORKDIR /app

# 设置容器的时区和时间
RUN apk --no-cache add tzdata
ENV TZ=Asia/Shanghai
COPY --from=builder /build/main /app/main
# 暴露端口
EXPOSE 8080
COPY . .
# 启动应用程序
CMD "/app/main"
