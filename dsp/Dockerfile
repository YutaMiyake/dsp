FROM asia.gcr.io/cyberagent-224/golang-build/golang-build:latest as build-env

ENV GOPATH /ws
WORKDIR /ws
RUN mkdir pkg && mkdir bin && mkdir src
COPY src src/.

RUN cd src/app/ && \
    rm -rf vendor && \
    dep ensure && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app -v -ldflags '-w -s' && \
    mv app /ws/app

FROM alpine
WORKDIR /root

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=build-env /ws/app /root/
COPY --from=build-env /ws/src/app/resources /root/resources

EXPOSE 8080
ADD run.sh /root/
CMD ["/root/run.sh"]
