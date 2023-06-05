FROM golang:1.20 as buildenv

WORKDIR /src
ADD . /src

# ENV APP_ENV=dev
# ENV APP_ENV=test
# ENV CLICKHOUSE_NAME=default            
# ENV CLICKHOUSE_HOST=127.0.0.1
# ENV CLICKHOUSE_PASSWORD=qwerty123
# ENV CLICKHOUSE_PORT=8123
# ENV CLICKHOUSE_USER=default
# RUN go mod tidy
# RUN go vet ./...
# RUN go test ./... -cover

RUN go build -o event ./cmd/main.go
RUN chmod +x event

# ////////////////////////////////////////////////////////////////
FROM ubuntu:latest
WORKDIR /usr/local/app

ENV APP_ENV=dev

COPY --from=buildenv /src/cmd/config-${APP_ENV}.toml ./config-${APP_ENV}.toml
COPY --from=buildenv /src/event .

EXPOSE 8080

CMD ["sh", "-c", "./event --cfg=config-dev.toml"]