# 使用 golang:1.19 镜像作为基础镜像
FROM golang:1.19 AS builder
# 设置工作目录
WORKDIR /app
# 将项目文件复制到容器内
COPY . .
# 设置代理环境变量
ARG GOPROXY="https://goproxy.cn,direct"

# 切换到子目录并使用交叉编译生成适用于 Linux 的可执行文件
RUN go build -o main main.go

# 最终的运行阶段
FROM alpine
# 从 builder 阶段复制可执行文件到最终的镜像中
COPY --from=builder /app/cmd/main /app/main

# 指定工作目录
WORKDIR /app

# 运行可执行文件
CMD ["./main"]
