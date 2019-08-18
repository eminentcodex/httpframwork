ARG APP_DIR=/go/src/httpframework
FROM golang:latest AS build-stage
ARG APP_DIR
RUN mkdir -p ${APP_DIR}
COPY . ${APP_DIR}
WORKDIR ${APP_DIR}
RUN export GOPATH=/go/src/ && export CGO_ENABLED=0 && export GOOS=linux && export GOMOD111=on &&\
 go mod tidy &&\
 go build -tags netgo -a -v -o httpframework_app .

FROM golang:latest AS production
ARG APP_DIR
WORKDIR ${APP_DIR}
COPY --from=build-stage ${APP_DIR}/httpframework_app ${APP_DIR}/server.crt ${APP_DIR}/server.key /go/bin/
COPY --from=build-stage ${APP_DIR}/configs/* /go/bin/configs/
RUN cat /go/bin/configs/config.yml
EXPOSE 8080/tcp
ENTRYPOINT ["/go/bin/httpframework_app"]
