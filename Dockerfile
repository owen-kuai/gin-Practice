# Build the manager binary
FROM 172.16.1.99/gold/golang:1.16 as builder
# FROM golang:1.15 as builder

WORKDIR /workspace
ENV GOPROXY="https://goproxy.io,https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct"
ENV GOSUMDB="off"

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

RUN go get -u github.com/swaggo/swag/cmd/swag && go mod download

# Copy the go source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on make build

# FROM gcr.io/distroless/static:nonroot
FROM centos8

WORKDIR /

RUN mkdir -p /etc/app/configs

COPY --from=builder /workspace/configs/config.yaml /etc/app/configs/config.yaml
COPY --from=builder /workspace/bin/app .

# USER nonroot:nonroot

# port
EXPOSE 8899

# entrypoint "app run -c  /etc/app/configs/config.yaml"
# ENTRYPOINT ["/app"]
