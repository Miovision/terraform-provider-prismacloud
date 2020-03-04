SHELL=/usr/bin/env bash

BUILD_DIR=build

provider: $(BUILD_DIR) $(BUILD_DIR)/terraform-provider-prismacloud-linux $(BUILD_DIR)/terraform-provider-prismacloud-mac

test:
	go test ./prismacloud_client/

fmt:
	go fmt . ./prismacloud/ ./prismacloud_client/

clean:
	rm -r $(BUILD_DIR)

$(BUILD_DIR):
	mkdir -p $@

$(BUILD_DIR)/terraform-provider-prismacloud-linux:
	GOOS=linux GOARCH=amd64 go build -o $@

$(BUILD_DIR)/terraform-provider-prismacloud-mac:
	GOOS=darwin GOARCH=amd64 go build -o $@
