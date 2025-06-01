# 构建阶段
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 复制 go mod 和 sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ar-backend .

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/ar-backend .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/docs ./docs

# 暴露端口
EXPOSE 3000

# 运行应用
CMD ["./ar-backend"] 