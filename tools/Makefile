.PHONY: proto gen install-protoc

PROTO_DIR := proto
GEN_DIR := gen/go
PROTOC := protoc

PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

proto: gen

gen: $(PROTO_FILES)
	@mkdir -p $(GEN_DIR)
	$(PROTOC) \
		--go_out=$(GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) \
		--go-grpc_opt=paths=source_relative \
		-I$(PROTO_DIR) \
		$(PROTO_FILES)
