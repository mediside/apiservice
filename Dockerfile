FROM golang:1.24-alpine3.22 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -C cmd/apiservice -o /apiservice

FROM alpine:3.22
WORKDIR /app
ENV CONFIG_PATH=./app/config/prod.yaml
COPY --from=builder /apiservice /app/apiservice
COPY config/ app/config/
EXPOSE 9042

ENTRYPOINT [ "/app/apiservice" ]