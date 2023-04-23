FROM golang:1.19

ENV GO111MODULE=on
RUN go install github.com/cucumber/godog/cmd/godog@v0.12.1
ADD https://github.com/dnnrly/wait-for/releases/download/v0.0.1/wait-for_0.0.1_linux_386.tar.gz wait-for.tar.gz
RUN gunzip wait-for.tar.gz && tar -xf wait-for.tar && mv wait-for /usr/local/bin

RUN mkdir /build
WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download

CMD godog