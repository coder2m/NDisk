FROM golang:1.14-alpine AS builder
RUN apk update && mkdir /app
WORKDIR /app
#ENV GOPROXY https://goproxy.cn
ENV GO111MODULE on
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 go build -o  dist/nfile -a -ldflags "-w -s -X github.com/coder2z/g-server/xapp.appName=ndisk_nfile\
                                                       -X github.com/coder2z/g-server/xapp.buildAppVersion=v1.0 \
                                                       -X github.com/coder2z/g-server/xapp.buildHost=`hostname -f`" ./cmd/nfile

FROM alpine:3.10

ARG REPO_URL
ARG BRANCH
ARG COMMIT_REF
LABEL repo-url=${REPO_URL}
LABEL branch=${BRANCH}
LABEL commit-ref=${COMMIT_REF}

RUN apk update \
  && apk upgrade \
  && apk add --no-cache \
  ca-certificates tzdata bash \
  && update-ca-certificates 2>/dev/null || true \
  && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
  && mkdir /app
WORKDIR /app
COPY --from=builder /app/dist/nfile /app/nfile

CMD ["/app/nfile"]