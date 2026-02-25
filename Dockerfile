# Build stage
FROM golang:1.23-alpine AS builder

# 设置构建参数
ARG VERSION=dev
ARG BUILD_TIME

# 设置工作目录
WORKDIR /build


# 复制源代码
COPY . .
RUN go mod tidy
RUN go mod vendor

# 构建应用
RUN go build \
    -ldflags "-s -w -X 'pkg.Version=${VERSION}' -X 'pkg.BuildTime=${BUILD_TIME}'" \
    -o /build/verge cmd/main.go

# Runtime stage
FROM alpine:3.19

# 安装必要的运行时依赖
#RUN apk add --no-cache ca-certificates tzdata

# 创建非 root 用户
RUN addgroup -g 1000 verge && \
    adduser -u 1000 -G verge -s /bin/sh -D verge

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /build/verge /app/verge

# 复制资源文件
#COPY --from=builder /build/res /app/res
#COPY --from=builder /build/platform/linux/start.sh /app/start.sh

# 设置权限
RUN chmod +x /app/verge /app/start.sh && \
    chown -R verge:verge /app

# 切换到非 root 用户
USER verge

# 暴露端口（根据应用需要调整）
# EXPOSE 8080

# 设置启动命令
CMD ["/app/verge"]