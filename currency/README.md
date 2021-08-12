# **Currency Service**
The currency service is a gRPC service which provides up to date exchange rates and currency conversion capabilities.

## **Installation for mac**
**Install protobuf compiler**

```shell
brew install protobuf
```

[gRPC-with-go](https://www.grpc.io/docs/languages/go/quickstart/)

```Go plugins``` **for the protocol compiler**

1. Install the protocol compiler plugins for Go using the following commands:

```shell
$ go install google.golang.org/protobuf/cmd/protoc-gen-go
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

2. Update your ```PATH``` so that the ```protoc``` compiler can find the plugins:

```shell
$ export PATH="$PATH:$(go env GOPATH)/bin"
```

**Adding module for project**

```shell
$ go get -u google.golang.org/grpc
```

**Then run the build command**

```shell
$ protoc -I $(PWD)/protos/ $(PWD)/protos/currency.proto \
    --go_out=$(PWD)/protos \
    --go-grpc_out=$(PWD)/protos
```

## **Testing**
**Installing the ```grpcurl```**

```shell
$ go install github.com/fullstorydev/grpcurl/cmd/grpcurl
```

**Run test**

```sh
$ grpcurl -plaintext localhost:9092 list
Currency

$ grpcurl -plaintext localhost:9092 list Currency
Currency.GetRate

$ grpcurl -plaintext -d '{"Base": "VND", "Destination": "USD"}' localhost:9092 Currency.GetRate
{
  "rate": 0.5
}
```
