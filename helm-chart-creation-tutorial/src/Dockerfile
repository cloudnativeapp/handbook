FROM golang:1.13-rc-alpine3.10 as builder
WORKDIR /app
COPY main.go .
RUN go build -o hello-world main.go

FROM alpine:3.10
WORKDIR /app
ARG PORT=80
COPY --from=builder /app/hello-world /app/hello-world
ENTRYPOINT ./hello-world
EXPOSE 80
