FROM golang:1.14-alpine AS builder

RUN apk add make upx
WORKDIR /app
COPY . /app
RUN make build
RUN upx dist/*

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/dist/api /app/api
COPY --from=builder /app/dist/notifier /app/notifier
COPY --from=builder /app/dist/session-cleaner /app/session-cleaner
COPY --from=builder /app/dist/updater /app/updater

COPY ./Makefile /app/Makefile
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
COPY ./sql /app/sql
