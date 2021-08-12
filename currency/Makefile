.PHONY: protos

protos:
	protoc -I $(PWD)/protos/ $(PWD)/protos/currency.proto \
    --go_out=$(PWD)/protos \
    --go-grpc_out=$(PWD)/protos