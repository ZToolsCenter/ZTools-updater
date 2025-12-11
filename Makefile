# 二进制文件名称
BINARY_NAME=ztools-updater

# 构建目录
BUILD_DIR=dist

.PHONY: all clean build-mac build-win

all: clean build-mac build-win

build-mac:
	@echo "Building for macOS (ARM64)..."
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/mac-arm64/$(BINARY_NAME) .

build-win:
	@echo "Building for Windows (AMD64)..."
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/win-amd64/$(BINARY_NAME).exe .

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
