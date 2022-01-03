FROM golang:1.16-alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
COPY . ./
RUN go mod download

RUN go build -o huawei-push-authorizator cmd/server/main.go 

RUN mv huawei-push-authorizator / && rm -Rf /build

EXPOSE 8077

CMD ["/huawei-push-authorizator"]