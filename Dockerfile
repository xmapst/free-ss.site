#docker build --rm --build-arg APP_ROOT=/go/src/free-ss -t free-ss:latest
FROM golang:1.16.0
ARG  APP_ROOT=/go/src/free-ss
WORKDIR ${APP_ROOT}
COPY ./ ${APP_ROOT}
ENV PATH=$GOPATH/bin:$PATH
RUN apt-get update \
  && apt-get install upx musl-dev -y

RUN export GO111MODULE=on \
  && go get -u github.com/swaggo/swag/cmd/swag \
  && swag init \
  && go mod tidy \
  && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main . \
  && strip --strip-unneeded main \
  && upx --lzma main

FROM alpine:latest
ARG  APP_ROOT=/go/src/free-ss
WORKDIR /app/
COPY --from=0 ${APP_ROOT}/main .
EXPOSE 80/tcp
ENTRYPOINT ["/app/main"]
