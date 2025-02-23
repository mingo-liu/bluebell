FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build

# 复制项目中的 go.mod 和 go.sum 文件并下载依赖信息
COPY go.mod . 
COPY go.sum . 
# 下载 go.mod 文件中列出的所有依赖项（模块）
RUN go mod download

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文件 bluebell
RUN go build -o bluebell .

###################
# 接下来创建一个小镜像
###################
FROM debian:bullseye-slim

# 复制脚本和配置文件到容器中
COPY ./wait-for.sh / 

# 拷贝配置文件
COPY ./conf /conf

# 从builder镜像中把静态文件拷贝到当前目录
COPY ./templates /templates
COPY ./static /static

# 从 builder 镜像中将可执行文件拷贝到当前目录
COPY --from=builder /build/bluebell / 

# 更新源并安装 netcat
RUN apt-get update \
    && apt-get install -y --no-install-recommends netcat \
    && chmod 755 wait-for.sh

# 设置运行的命令（可根据需要调整）
# ENTRYPOINT ["/bluebell", "conf/config.ini"]
