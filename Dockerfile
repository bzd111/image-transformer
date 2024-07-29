FROM golang:1.22.1-alpine3.18 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN  go build -o webhook main.go

FROM alpine:3.18.0

ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/webhook /app/webhook

EXPOSE 8443

ENTRYPOINT  ["./webhook"]
