from sysdig/sysdig
maintainer james.xiong@daocloud.io

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR /gopath/app
ENV GOPATH /gopath/app
ADD . /gopath/app/
RUN go install dockermetrics
RUN cp /gopath/app/src/templates/* /gopath/app/bin/
RUN rm -fr /gopath/app/src
ENV TZ Asia/Shanghai
EXPOSE 8080
WORKDIR /gopath/app/bin
#entrypoint ["./dockermetrics"]
cmd ["./dockermetrics"]